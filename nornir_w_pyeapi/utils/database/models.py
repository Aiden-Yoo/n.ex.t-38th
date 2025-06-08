from nornir.core.task import Result, Task
from sqlalchemy import Column, Integer, String, DateTime, Text, create_engine
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime, UTC
import json
import logging

logger = logging.getLogger("nornir")

Base = declarative_base()

class TaskResult(Base):
    __abstract__ = True

    id = Column(Integer, primary_key=True)
    hostname = Column(String(255), nullable=False)
    task_name = Column(String(100), nullable=False)
    output = Column(Text)
    status = Column(String(20), nullable=False) 
    error = Column(Text) 
    created_at = Column(DateTime, default=lambda: datetime.now(UTC))

    def __repr__(self):
        return f"<TaskResult(hostname='{self.hostname}', task_name='{self.task_name}', status='{self.status}')>"

def get_task_result_class(task_name: str):
    class_name = f"{task_name.capitalize()}Result"
    return type(
        class_name,
        (TaskResult,),
        {
            '__tablename__': task_name,
            '__table_args__': {'extend_existing': True}
        }
    )

def save_task_result(task: Task, result: Result):
    from .db_conn import db
    
    with db.get_session() as session:
        if isinstance(result, dict):
            status = 'success'
            output = json.dumps(result, ensure_ascii=False)
            error = None
        else:
            status = 'success' if not result.failed else 'failed'
            output = json.dumps(result.result, ensure_ascii=False) if not result.failed else None
            error = str(result.exception) if result.failed else None

        TaskResultClass = get_task_result_class(task.name)
        
        Base.metadata.create_all(db.engine, tables=[TaskResultClass.__table__])
        
        task_result = TaskResultClass(
            hostname=task.host.hostname,
            task_name=task.name,
            output=output,
            status=status,
            error=error
        )
        session.add(task_result)
        logger.info(f"Saved result for {task.host.hostname} - {task.name}") 