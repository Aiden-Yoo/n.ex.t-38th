from typing import Dict
from nornir.core.task import Result, Task
from utils.pyeapi_conn import get_eapi_connection

def config_intf(task: Task, interfaces: Dict[str, Dict[str, str]]) -> Result:
    try:
        connection = get_eapi_connection(task)
        commands = []
        
        commands.append('configure')
        
        for interface, config in interfaces.items():
            commands.append(f'interface {interface}')
            for key, value in config.items():
                if value:
                    commands.append(f'{key} {value}')
                else:
                    commands.append(key)
            commands.append('exit') 

        result = connection.execute(commands, encoding='json')
        return Result(host=task.host, result=result)
    except Exception as e:
        return Result(host=task.host, failed=True, exception=e)