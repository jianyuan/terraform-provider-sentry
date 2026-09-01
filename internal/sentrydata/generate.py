import ast
import asyncio
import pathlib
import zoneinfo
from collections import OrderedDict
from collections.abc import Mapping
from typing import Any, Generic, NamedTuple, TypeVar

import httpx2
from mako.template import Template

REPO = "getsentry/sentry"
BRANCH = "master"

ROOT = pathlib.Path(__file__).parent.resolve()
OUTPUT_FILE = (ROOT / "sentrydata.go").resolve()
TEMPLATE_FILE = (ROOT / "sentrydata.go.mako").resolve()

TEMPLATE = Template(filename=str(TEMPLATE_FILE))


class FileData(NamedTuple):
    github_url: str
    tree: ast.Module


DataT = TypeVar("DataT", list[str], OrderedDict[str, str])


class Variable(NamedTuple, Generic[DataT]):
    github_url: str
    data: DataT

    @staticmethod
    def to_go_type(value: Any) -> str:
        match value:
            case str():
                return "string"
            case int():
                return "int64"
            case bool():
                return "bool"
            case _:
                raise ValueError(f"Value type {type(value)} is not supported")

    @property
    def go_key_type(self) -> str:
        assert isinstance(self.data, Mapping), f"Data {self.data} is not a Mapping"
        go_types = {self.to_go_type(v) for v in self.data}
        assert len(go_types) == 1, f"Key types must be the same, got {go_types}"
        return next(iter(go_types))

    @property
    def go_value_type(self) -> str:
        values = self.data.values() if isinstance(self.data, Mapping) else self.data
        go_types = {self.to_go_type(v) for v in values}
        assert len(go_types) == 1, f"Value types must be the same, got {go_types}"
        return next(iter(go_types))


async def get_file_data(client: httpx2.AsyncClient, path: str) -> FileData:
    url = f"https://raw.githubusercontent.com/{REPO}/refs/heads/{BRANCH}/{path}"
    r = await client.get(url)
    r.raise_for_status()
    tree = ast.parse(r.text)
    return FileData(
        github_url=f"https://github.com/{REPO}/blob/{BRANCH}/{path}",
        tree=tree,
    )


async def parse_constants(client: httpx2.AsyncClient) -> dict[str, Variable[Any]]:
    import logging

    data = await get_file_data(client, "src/sentry/constants.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.AnnAssign(
                target=ast.Name(id="LOG_LEVELS"),
                value=ast.Dict(keys=keys, values=values),
            ):
                out["LogLevels"] = Variable(
                    github_url=data.github_url,
                    data=[],
                )
                out["LogLevelNameToId"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                out["LogLevelIdToName"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                for key, value in zip(keys, values):
                    assert isinstance(key, ast.Attribute)
                    assert isinstance(value, ast.Constant)
                    log_level_id = str(getattr(logging, key.attr))
                    out["LogLevels"].data.append(value.value)
                    out["LogLevelNameToId"].data[value.value] = log_level_id
                    out["LogLevelIdToName"].data[log_level_id] = value.value
            case _:
                pass
    return out


async def parse_issues_grouptype(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/issues/grouptype.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name="GroupCategory", body=body):
                out["IssueGroupCategories"] = Variable(
                    github_url=data.github_url,
                    data=[],
                )
                out["IssueGroupCategoryNameToId"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                out["IssueGroupCategoryIdToName"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                for node in body:
                    match node:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Constant(value=value),
                        ) if id.isupper():
                            name = id.replace("_", " ").title().replace(" ", "_")
                            out["IssueGroupCategories"].data.append(name)
                            out["IssueGroupCategoryNameToId"].data[name] = str(value)
                            out["IssueGroupCategoryIdToName"].data[str(value)] = name
                        case _:
                            pass
            case _:
                pass
    return out


async def parse_rules_conditions_event_attribute(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/rules/conditions/event_attribute.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.AnnAssign(
                target=ast.Name(id="ATTR_CHOICES"),
                value=ast.Dict(keys=keys),
            ):
                out["EventAttributes"] = Variable(
                    github_url=data.github_url,
                    data=[],
                )
                for key in keys:
                    assert isinstance(key, ast.Constant)
                    out["EventAttributes"].data.append(key.value)
            case _:
                pass
    return out


async def parse_rules_match(client: httpx2.AsyncClient) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/rules/match.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name="MatchType", body=body):
                out["MatchTypes"] = Variable(
                    github_url=data.github_url,
                    data=[],
                )
                out["MatchTypeIds"] = Variable(
                    github_url=data.github_url,
                    data=[],
                )
                out["MatchTypeNameToId"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                out["MatchTypeIdToName"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                for node in body:
                    match node:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Constant(value=value),
                        ) if id.isupper():
                            out["MatchTypes"].data.append(id)
                            out["MatchTypeIds"].data.append(value)
                            out["MatchTypeNameToId"].data[id] = value
                            out["MatchTypeIdToName"].data[value] = id
                        case _:
                            pass
            case ast.Assign(
                targets=[ast.Name(id="LEVEL_MATCH_CHOICES")],
                value=ast.Dict(keys=keys),
            ):
                out["LevelMatchTypes"] = Variable(
                    github_url=data.github_url,
                    data=[],
                )
                for key in keys:
                    assert isinstance(key, ast.Attribute)
                    out["LevelMatchTypes"].data.append(key.attr)
            case _:
                pass
    return out


async def parse_models_dashboard_widget(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
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
                            ) if id.isupper():
                                out.append(value)
                            case _:
                                pass
                case _:
                    pass
        return out

    data = await get_file_data(client, "src/sentry/models/dashboard_widget.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name=name):
                types = extract_types(node)
                if types:
                    out[name] = Variable(
                        github_url=data.github_url,
                        data=types,
                    )
            case _:
                pass
    return out


async def parse_models_project(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/models/project.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.Assign(
                targets=[ast.Name(id="GETTING_STARTED_DOCS_PLATFORMS")],
                value=ast.List(elts=elts),
            ):
                out["Platforms"] = Variable(
                    github_url=data.github_url,
                    data=["other"],
                )
                for elt in elts:
                    assert isinstance(elt, ast.Constant)
                    out["Platforms"].data.append(elt.value)
            case _:
                pass
    return out


async def parse_intervals(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/monitors/validators.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.Assign(
                targets=[ast.Name(id="INTERVAL_NAMES")],
                value=ast.Tuple(elts=elts),
            ):
                result: list[str] = []
                for elt in elts:
                    match elt:
                        case ast.Constant(value=value):
                            result.append(value)
                        case _:
                            pass
                out["Intervals"] = Variable(github_url=data.github_url, data=result)
            case _:
                pass
    return out


async def get_timezones() -> dict[str, Variable[Any]]:
    timezones = frozenset(zoneinfo.available_timezones() - {"Factory", "localtime"})
    return {
        "Timezones": Variable(
            github_url=None,
            data=sorted(timezones),
        )
    }


async def parse_alert_rule_models(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/incidents/models/alert_rule.py")
    out: dict[str, Variable[Any]] = {
        "AlertRuleDetectionTypes": Variable(github_url=data.github_url, data=[]),
        "AlertRuleSensitivities": Variable(github_url=data.github_url, data=[]),
        "AlertRuleThresholdTypes": Variable(github_url=data.github_url, data=[]),
        "AlertRuleThresholdTypeNameToId": Variable(github_url=data.github_url, data={}),
        "AlertRuleThresholdTypeIdToName": Variable(github_url=data.github_url, data={}),
    }
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name="AlertRuleDetectionType", body=elts):
                for elt in elts:
                    match elt:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Tuple(elts=[ast.Constant(value=value), _]),
                        ) if id.isupper():
                            out["AlertRuleDetectionTypes"].data.append(value)
                        case _:
                            pass
            case ast.ClassDef(name="AlertRuleSensitivity", body=elts):
                for elt in elts:
                    match elt:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Tuple(elts=[ast.Constant(value=value), _]),
                        ) if id.isupper():
                            out["AlertRuleSensitivities"].data.append(value)
                        case _:
                            pass
            case ast.ClassDef(name="AlertRuleThresholdType", body=elts):
                for elt in elts:
                    match elt:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Constant(value=value),
                        ) if id.isupper():
                            out["AlertRuleThresholdTypes"].data.append(id.lower())
                            out["AlertRuleThresholdTypeNameToId"].data[id.lower()] = (
                                value
                            )
                            out["AlertRuleThresholdTypeIdToName"].data[value] = (
                                id.lower()
                            )
                        case _:
                            pass
            case _:
                pass
    return out


async def parse_uptime_models(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/uptime/models.py")
    out: dict[str, Variable[Any]] = {
        "UptimeSubscriptionSupportedHttpMethods": Variable(
            github_url=data.github_url,
            data=[],
        ),
        "UptimeSubscriptionIntervalSeconds": Variable(
            github_url=data.github_url,
            data=[],
        ),
    }
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(
                name="UptimeSubscription",
                body=body,
            ):
                for body_node in body:
                    match body_node:
                        case ast.ClassDef(name="SupportedHTTPMethods", body=elts):
                            for elt in elts:
                                match elt:
                                    case ast.Assign(
                                        value=ast.Tuple(
                                            elts=[ast.Constant(value=value), _]
                                        )
                                    ):
                                        out[
                                            "UptimeSubscriptionSupportedHttpMethods"
                                        ].data.append(value)
                                    case _:
                                        pass
                        case ast.ClassDef(name="IntervalSeconds", body=elts):
                            for elt in elts:
                                match elt:
                                    case ast.Assign(
                                        value=ast.Tuple(
                                            elts=[ast.Constant(value=value), _]
                                        )
                                    ):
                                        out[
                                            "UptimeSubscriptionIntervalSeconds"
                                        ].data.append(value)
                                    case _:
                                        pass
                        case _:
                            pass
            case _:
                pass
    return out


async def parse_uptime_types(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/uptime/types.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(
                name="UptimeMonitorMode",
                body=elts,
            ):
                out["UptimeMonitorModes"] = Variable(
                    github_url=data.github_url,
                    data=[],
                )
                out["UptimeMonitorModeNameToId"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                out["UptimeMonitorModeIdToName"] = Variable(
                    github_url=data.github_url,
                    data=OrderedDict(),
                )
                for elt in elts:
                    match elt:
                        case ast.Assign(
                            targets=[ast.Name(id=id)],
                            value=ast.Constant(value=value),
                        ) if id.isupper():
                            out["UptimeMonitorModes"].data.append(id)
                            out["UptimeMonitorModeNameToId"].data[id] = value
                            out["UptimeMonitorModeIdToName"].data[value] = id
                        case _:
                            pass
            case _:
                pass
    return out


async def parse_snuba_models(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/snuba/models.py")
    out: dict[str, Variable[Any]] = {
        "SnubaQueryTypes": Variable(github_url=data.github_url, data=[]),
        "SnubaQueryTypeNameToId": Variable(github_url=data.github_url, data={}),
        "SnubaQueryTypeIdToName": Variable(github_url=data.github_url, data={}),
        "SnubaExtrapolationModes": Variable(github_url=data.github_url, data=[]),
        "SnubaQueryEventTypes": Variable(github_url=data.github_url, data=[]),
        "SnubaQueryEventTypeNameToId": Variable(github_url=data.github_url, data={}),
        "SnubaQueryEventTypeIdToName": Variable(github_url=data.github_url, data={}),
    }
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(name="SnubaQuery", body=elts):
                for elt in elts:
                    match elt:
                        case ast.ClassDef(name="Type", body=type_elts):
                            for elt in type_elts:
                                match elt:
                                    case ast.Assign(
                                        targets=[ast.Name(id=id)],
                                        value=ast.Constant(value=value),
                                    ) if id.isupper():
                                        out["SnubaQueryTypes"].data.append(id.lower())
                                        out["SnubaQueryTypeNameToId"].data[
                                            id.lower()
                                        ] = value
                                        out["SnubaQueryTypeIdToName"].data[value] = (
                                            id.lower()
                                        )
                        case _:
                            pass
            case ast.ClassDef(name="SnubaQueryEventType", body=elts):
                for elt in elts:
                    match elt:
                        case ast.ClassDef(name="EventType", body=type_elts):
                            for elt in type_elts:
                                match elt:
                                    case ast.Assign(
                                        targets=[ast.Name(id=id)],
                                        value=ast.Constant(value=value),
                                    ) if id.isupper():
                                        out["SnubaQueryEventTypes"].data.append(
                                            id.lower()
                                        )
                                        out["SnubaQueryEventTypeNameToId"].data[
                                            id.lower()
                                        ] = value
                                        out["SnubaQueryEventTypeIdToName"].data[
                                            value
                                        ] = id.lower()
                        case _:
                            pass
            case ast.ClassDef(name="ExtrapolationMode", body=elts):
                for elt in elts:
                    match elt:
                        case ast.Assign(targets=[ast.Name(id=id)]) if id.isupper():
                            out["SnubaExtrapolationModes"].data.append(id.lower())
                        case _:
                            pass
            case _:
                pass
    return out


async def parse_snuba_datasets(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/snuba/dataset.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(
                name="Dataset",
                body=elts,
            ):
                result: list[str] = []
                for elt in elts:
                    match elt:
                        case ast.Assign(value=ast.Constant(value=value)):
                            result.append(value)
                        case _:
                            pass
                out["SnubaDatasets"] = Variable(github_url=data.github_url, data=result)
            case _:
                pass
    return out


async def parse_data_condition_types(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(
        client, "src/sentry/workflow_engine/models/data_condition.py"
    )
    out: dict[str, Variable[Any]] = {
        "DataConditionTypes": Variable(github_url=data.github_url, data=[]),
    }
    condition_types: dict[str, str] = {}

    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(
                name="Condition",
                body=elts,
            ):
                for elt in elts:
                    match elt:
                        case ast.Assign(
                            targets=[ast.Name(id=key)],
                            value=ast.Constant(value=value),
                        ):
                            out["DataConditionTypes"].data.append(value)
                            condition_types[key] = value
                        case _:
                            pass
            case _:
                pass
    return out


async def parse_data_condition_group_types(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(
        client, "src/sentry/workflow_engine/models/data_condition_group.py"
    )
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.ClassDef(
                name="DataConditionGroup",
                body=body,
            ):
                result: list[str] = []
                for body_node in body:
                    match body_node:
                        case ast.ClassDef(name="Type", body=elts):
                            for elt in elts:
                                match elt:
                                    case ast.Assign(value=ast.Constant(value=value)):
                                        result.append(value)
                                    case _:
                                        pass
                        case _:
                            pass
                out["DataConditionGroupTypes"] = Variable(
                    github_url=data.github_url, data=result
                )
            case _:
                pass
    return out


async def parse_event_frequency(
    client: httpx2.AsyncClient,
) -> dict[str, Variable[Any]]:
    data = await get_file_data(client, "src/sentry/rules/conditions/event_frequency.py")
    out: dict[str, Variable[Any]] = {}
    for node in ast.walk(data.tree):
        match node:
            case ast.AnnAssign(
                target=ast.Name(id="STANDARD_INTERVALS"),
                value=ast.Dict(keys=keys),
            ):
                result_intervals: list[str] = []
                for key in keys:
                    assert isinstance(key, ast.Constant)
                    result_intervals.append(key.value)
                out["EventFrequencyStandardIntervals"] = Variable(
                    github_url=data.github_url, data=result_intervals
                )
            case _:
                pass
    return out


async def generate_and_format_code(variables: dict[str, Variable[Any]]) -> bytes:
    unformatted_code = TEMPLATE.render(variables=variables)

    gofmt_proc = await asyncio.create_subprocess_exec(
        "gofmt",
        stdin=asyncio.subprocess.PIPE,
        stdout=asyncio.subprocess.PIPE,
        stderr=asyncio.subprocess.PIPE,
    )

    formatted_code, stderr = await gofmt_proc.communicate(
        input=unformatted_code.encode()
    )

    if gofmt_proc.returncode != 0:
        raise RuntimeError(f"gofmt failed:\n{stderr.decode()}")

    return formatted_code


async def main() -> None:
    async with httpx2.AsyncClient(
        headers={"User-Agent": "terraform-provider-sentry-data-generator/dev"}
    ) as client:
        results = await asyncio.gather(
            parse_constants(client=client),
            parse_issues_grouptype(client=client),
            parse_rules_conditions_event_attribute(client=client),
            parse_rules_match(client=client),
            parse_models_dashboard_widget(client=client),
            parse_models_project(client=client),
            parse_intervals(client=client),
            get_timezones(),
            parse_alert_rule_models(client=client),
            parse_uptime_models(client=client),
            parse_uptime_types(client=client),
            parse_snuba_models(client=client),
            parse_snuba_datasets(client=client),
            parse_data_condition_types(client=client),
            parse_data_condition_group_types(client=client),
            parse_event_frequency(client=client),
        )

    variables: OrderedDict[str, Variable[Any]] = OrderedDict()
    for result in results:
        variables.update(result)

    code = await generate_and_format_code(variables=variables)
    OUTPUT_FILE.write_bytes(code)


if __name__ == "__main__":
    asyncio.run(main())
