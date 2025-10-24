from __future__ import annotations

import ast
import pathlib
import subprocess
from typing import Any, Generic, NamedTuple, OrderedDict, TypeGuard, TypeVar

import httpx
import jinja2

REPO = "getsentry/sentry"
BRANCH = "master"
TEMPLATE = """package sentrydata

{% for key, value in result.items() %}
// {{ value.github_url }}
var {{ key }}{{ ' = ' }}{% if value.result|is_list %}
[]string{
{% for v in value.result %}
    "{{ v }}",
{% endfor %}
}
{% else %}
map[string]string{
{% for k, v in value.result.items() %}
    "{{ k }}": "{{ v }}",
{% endfor %}
}
{% endif %}

{% endfor %}
"""


def get_jinja2_env() -> jinja2.Environment:
    env = jinja2.Environment(
        trim_blocks=True,
        lstrip_blocks=True,
    )

    def is_list(value: Any) -> TypeGuard[list[Any]]:
        return isinstance(value, list)

    env.filters["is_list"] = is_list

    return env


def get_text(path: str) -> str:
    r = httpx.get(
        f"https://raw.githubusercontent.com/{REPO}/refs/heads/{BRANCH}/{path}"
    )
    r.raise_for_status()
    return r.text


class FileData(NamedTuple):
    github_url: str
    tree: ast.Module


def get_file_data(path: str) -> FileData:
    return FileData(
        github_url=f"https://github.com/{REPO}/blob/{BRANCH}/{path}",
        tree=ast.parse(get_text(path)),
    )


ResultT = TypeVar("ResultT", list[str], OrderedDict[str, str])


class ResultData(NamedTuple, Generic[ResultT]):
    github_url: str
    result: ResultT


def parse_constants() -> dict[str, ResultData[Any]]:
    import logging

    data = get_file_data("src/sentry/constants.py")
    out: dict[str, ResultData[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.Assign(
                targets=[ast.Name(id="LOG_LEVELS")],
                value=ast.Dict(keys=keys, values=values),
            ):
                out["LogLevels"] = ResultData(
                    github_url=data.github_url,
                    result=[],
                )
                out["LogLevelNameToId"] = ResultData(
                    github_url=data.github_url,
                    result=OrderedDict(),
                )
                out["LogLevelIdToName"] = ResultData(
                    github_url=data.github_url,
                    result=OrderedDict(),
                )
                for key, value in zip(keys, values):
                    assert isinstance(key, ast.Attribute)
                    assert isinstance(value, ast.Constant)
                    log_level_id = str(getattr(logging, key.attr))
                    out["LogLevels"].result.append(value.value)
                    out["LogLevelNameToId"].result[value.value] = log_level_id
                    out["LogLevelIdToName"].result[log_level_id] = value.value
            case _:
                pass
    return out


def parse_issues_grouptype() -> dict[str, ResultData[Any]]:
    data = get_file_data("src/sentry/issues/grouptype.py")
    out: dict[str, ResultData[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name="GroupCategory", body=body):
                out["IssueGroupCategories"] = ResultData(
                    github_url=data.github_url,
                    result=[],
                )
                out["IssueGroupCategoryNameToId"] = ResultData(
                    github_url=data.github_url,
                    result=OrderedDict(),
                )
                out["IssueGroupCategoryIdToName"] = ResultData(
                    github_url=data.github_url,
                    result=OrderedDict(),
                )
                for node in body:
                    match node:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Constant(value=value),
                        ) if (
                            id.upper() == id
                        ):
                            name = id.replace("_", " ").title().replace(" ", "_")
                            out["IssueGroupCategories"].result.append(name)
                            out["IssueGroupCategoryNameToId"].result[name] = str(value)
                            out["IssueGroupCategoryIdToName"].result[str(value)] = name
                        case _:
                            pass
            case _:
                pass
    return out


def parse_rules_conditions_event_attribute() -> dict[str, ResultData[Any]]:
    data = get_file_data("src/sentry/rules/conditions/event_attribute.py")
    out: dict[str, ResultData[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.AnnAssign(
                target=ast.Name(id="ATTR_CHOICES"),
                value=ast.Dict(keys=keys),
            ):
                out["EventAttributes"] = ResultData(
                    github_url=data.github_url,
                    result=[],
                )
                for key in keys:
                    assert isinstance(key, ast.Constant)
                    out["EventAttributes"].result.append(key.value)
            case _:
                pass
    return out


def parse_rules_match() -> dict[str, ResultData[Any]]:
    data = get_file_data(path="src/sentry/rules/match.py")
    out: dict[str, ResultData[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name="MatchType", body=body):
                out["MatchTypes"] = ResultData(
                    github_url=data.github_url,
                    result=[],
                )
                out["MatchTypeNameToId"] = ResultData(
                    github_url=data.github_url,
                    result=OrderedDict(),
                )
                out["MatchTypeIdToName"] = ResultData(
                    github_url=data.github_url,
                    result=OrderedDict(),
                )
                for node in body:
                    match node:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Constant(value=value),
                        ) if (
                            id.upper() == id
                        ):
                            out["MatchTypes"].result.append(id)
                            out["MatchTypeNameToId"].result[id] = value
                            out["MatchTypeIdToName"].result[value] = id
                        case _:
                            pass
            case ast.Assign(
                targets=[ast.Name(id="LEVEL_MATCH_CHOICES")],
                value=ast.Dict(keys=keys),
            ):
                out["LevelMatchTypes"] = ResultData(
                    github_url=data.github_url,
                    result=[],
                )
                for key in keys:
                    assert isinstance(key, ast.Attribute)
                    out["LevelMatchTypes"].result.append(key.attr)
            case _:
                pass
    return out


def parse_models_dashboard_widget() -> dict[str, ResultData[Any]]:
    def extract_types(classdef: ast.ClassDef) -> list[str]:
        out: list[str] = []
        for node in classdef.body:
            match node:
                case ast.Assign(
                    targets=[ast.Name(id="TYPES")],
                    value=ast.List(elts=elts),
                ):
                    for elt in elts:
                        match elt:
                            case ast.Tuple(
                                elts=[ast.Name(id=id), ast.Constant(value=value)]
                            ) if (id.upper() == id):
                                out.append(value)
                            case _:
                                pass
                case _:
                    pass
        return out

    data = get_file_data("src/sentry/models/dashboard_widget.py")
    out: dict[str, ResultData[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name=name):
                types = extract_types(node)
                if types:
                    out[name] = ResultData(
                        github_url=data.github_url,
                        result=types,
                    )
            case _:
                pass
    return out


def parse_models_project() -> dict[str, ResultData[Any]]:
    data = get_file_data("src/sentry/models/project.py")
    out: dict[str, ResultData[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.Assign(
                targets=[ast.Name(id="GETTING_STARTED_DOCS_PLATFORMS")],
                value=ast.List(elts=elts),
            ):
                out["Platforms"] = ResultData(
                    github_url=data.github_url,
                    result=["other"],
                )
                for elt in elts:
                    assert isinstance(elt, ast.Constant)
                    out["Platforms"].result.append(elt.value)
            case _:
                pass
    return out


def main() -> None:
    result: OrderedDict[str, ResultData[Any]] = OrderedDict()
    result.update(parse_constants())
    result.update(parse_issues_grouptype())
    result.update(parse_rules_conditions_event_attribute())
    result.update(parse_rules_match())
    result.update(parse_models_dashboard_widget())
    result.update(parse_models_project())

    env = get_jinja2_env()
    template = env.from_string(TEMPLATE)
    output = (pathlib.Path(__file__).parent / "sentrydata.go").resolve()

    with output.open("w") as f:
        f.write(template.render(result=result))

    subprocess.run(["gofmt", "-w", output])


if __name__ == "__main__":
    main()
