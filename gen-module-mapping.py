from pathlib import Path
import json
import re

module_path = Path("./FMD2/lua/modules")

var_names = ["m.ID", "m.Name", "m.RootURL"]

def getVariables(module_data):
    variables = {}
    var_pattern = re.compile(r'\s*(?P<name>.*?)\s*=\s*(?:["\'])(?P<value>.*?)(?:["\'])')
    matches = var_pattern.finditer(module_data)
    for m in matches:
        group = m.groupdict()
        variables[group["name"]] = group["value"]
    try:
        data = {key: variables[key] for key in var_names}
    except Exception as e:
        return None
    return data

def getVariablesViaWebsiteModule(module_data):
    addcall_pattern = re.compile(r'\s*AddWebsiteModule\((?:["\'])(?P<id>.*?)(?:["\']),\s*(?:["\'])(?P<name>.*?)(?:["\']),\s*(?:["\'])(?P<rootURL>.*?)(?:["\'])')
    matches = addcall_pattern.finditer(module_data)
    result = []
    for m in matches:
        group = m.groupdict()
        result.append({"m.ID": group["id"], "m.Name": group["name"], "m.RootURL": group["rootURL"]})
    return result
        

result = []
for file in module_path.glob("**/*.lua"):
    print(file)
    module_data = file.open('r', encoding='ISO-8859-1').read()
    
    data = getVariables(module_data)
    if not data:
        data = getVariablesViaWebsiteModule(module_data)
        result.extend(data)
    else:
        result.append(data)

with Path("mapping.json").open('w+') as f:
    json.dump(result, f, indent=4)
