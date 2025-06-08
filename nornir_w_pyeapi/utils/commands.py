from typing import List
from nornir.core.task import Result, Task
from utils.pyeapi_conn import get_eapi_connection
import logging

logger = logging.getLogger("nornir")

def execute_all(task: Task, commands: List[str]) -> Result:
    try:
        connection = get_eapi_connection(task)
        logger.info(f"Executing commands on {task.host.hostname}: {commands}")
        result = connection.execute(commands)
        logger.info(f"Successfully executed commands on {task.host.hostname}")
        return Result(host=task.host, result=result)
    except Exception as e:
        logger.error(f"Failed to execute commands on {task.host.hostname}: {str(e)}")
        return Result(host=task.host, failed=True, exception=e)


