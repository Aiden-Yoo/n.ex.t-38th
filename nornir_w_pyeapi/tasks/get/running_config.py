from nornir.core.task import Result, Task
from utils.pyeapi_conn import get_eapi_connection
from utils.commands import execute_all
import logging

logger = logging.getLogger("nornir")

def get_running_config_raw(task: Task) -> Result:
    return execute_all(task, ['show running-config']) 

def get_running_config(task: Task) -> Result:
    command = 'show running-config'
    try:
        connection = get_eapi_connection(task)
        logger.info(f"Executing commands on {task.host.hostname}: {command}")
        output = connection.execute(command)
        logger.info(f"Successfully executed commands on {task.host.hostname}")
        return Result(host=task.host, result=output)
    except Exception as e:
        logger.error(f"Failed to execute commands on {task.host.hostname}: {str(e)}")
        return Result(host=task.host, failed=True, exception=e)
