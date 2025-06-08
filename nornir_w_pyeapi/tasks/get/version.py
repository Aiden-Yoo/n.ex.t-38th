from nornir.core.task import Result, Task
from utils.pyeapi_conn import get_eapi_connection
from utils.commands import execute_all
import humanize
import logging
from utils.database.models import save_task_result

logger = logging.getLogger("nornir")

def get_version_raw(task: Task) -> Result:
    return execute_all(task, ['show version'])

def get_version(task: Task, save_result: bool = False) -> Result:
    command = 'show version'
    try:
        connection = get_eapi_connection(task)
        logger.info(f"Executing commands on {task.host.hostname}: {command}")
        output = connection.execute(command)
        logger.info(f"Successfully executed commands on {task.host.hostname}")

        result = output.get('result', {})[0]

        version = result.get('version', {})
        model_name = result.get('modelName', {})
        serial_number = result.get('serialNumber', {})
        uptime = result.get('uptime', {})
        mem_total = result.get('memTotal', {})
        mem_free = result.get('memFree', {})
        mem_used = (mem_total - mem_free)/mem_total*100

        custom_output = {
            'version': version,
            'modelName': model_name,
            'serialNumber': serial_number,
            'uptime': humanize.precisedelta(uptime),
            'memUsed': "{:.2f}%".format(mem_used),
        }

        if save_result:
            save_task_result(task, custom_output)

        return Result(host=task.host, result=custom_output)
    except Exception as e:
        logger.error(f"Failed to execute commands on {task.host.hostname}: {str(e)}")
        return Result(host=task.host, failed=True, exception=e)
