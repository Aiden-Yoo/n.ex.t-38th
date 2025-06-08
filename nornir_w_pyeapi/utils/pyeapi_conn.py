from nornir.core.task import Task
import pyeapi
import logging

logger = logging.getLogger("nornir")

def get_eapi_connection(task: Task) -> pyeapi.client.Node:
    try:
        logger.info(f"Connecting to {task.host.hostname}")
        connection = pyeapi.connect(
            host=task.host.hostname,
            username=task.host.username,
            password=task.host.password,
            transport='https',
            port=443,
            verify=False
        )
        logger.info(f"Successfully connected to {task.host.hostname}")
        return connection
    except Exception as e:
        logger.error(f"Failed to connect to {task.host.hostname}: {str(e)}")
        raise