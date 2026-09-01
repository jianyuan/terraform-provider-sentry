<%! import json; from collections.abc import Mapping %>
  package sentrydata

  % for key, value in variables.items():
  % if value.github_url:
  // ${value.github_url}
  % endif
  % if isinstance(value.data, Mapping):
  var ${key} = map[${value.go_key_type}]${value.go_value_type}{
  % for k, v in value.data.items():
  ${json.dumps(k)}: ${json.dumps(v)},
  % endfor
  }
  % else:
  var ${key} = []${value.go_value_type}{
  % for v in value.data:
  ${json.dumps(v)},
  % endfor
  }
  % endif

  % endfor