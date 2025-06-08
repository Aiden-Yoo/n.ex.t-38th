from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, Session
from contextlib import contextmanager
from typing import Generator
import logging

logger = logging.getLogger("nornir")

class DatabaseConnection:
    def __init__(self, db_url: str = 'sqlite:///nornir_results.db'):
        self.engine = create_engine(db_url)
        self.Session = sessionmaker(bind=self.engine)
    
    def init_db(self):
        """데이터베이스 초기화"""
        from .models import Base
        Base.metadata.create_all(self.engine)
    
    @contextmanager
    def get_session(self) -> Generator[Session, None, None]:
        session = self.Session()
        try:
            yield session
            session.commit()
        except Exception as e:
            session.rollback()
            logger.error(f"Database error: {str(e)}")
            raise
        finally:
            session.close()

db = DatabaseConnection() 