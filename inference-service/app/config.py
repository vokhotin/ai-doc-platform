from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    host: str = '0.0.0.0'
    port: int = 8001

def get_settings() -> Settings:
    return Settings()