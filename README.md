# Descarga Concurrente de Imágenes con Go y React

Este proyecto es una solución al desafío de desarrollar una aplicación en Go que descargue imágenes de múltiples URLs simultáneamente, permitiendo al usuario configurar el número deseado de descargas concurrentes. Utiliza Go para el backend y React para el frontend, implementando principios REST y mejores prácticas de desarrollo moderno.

## Características Principales

- **Descargas Concurrentes Configurables**: Permite al usuario establecer y modificar el límite de descargas concurrentes.
- **Indicador de Progreso**: Ofrece información detallada sobre el progreso de cada descarga.
- **Manejo de Errores**: Asegura respuestas HTTP adecuadas en caso de fallos en las descargas.
- **Interfaz de Usuario con React**: Proporciona una experiencia de usuario fluida y responsiva para la interacción con la aplicación.

## Tecnologías Utilizadas

- **Backend**: Go
- **Frontend**: React
- **Mensajería en Tiempo Real**: Redis/PubSub
- **Persistencia de Datos**: MariaDB
- **Contenerización**: Docker

## Instrucciones de Instalación

1. **Clonar el Repositorio**

```bash
git clone https://github.com/gouh/concurrent-image-fetcher.git
cd concurrent-image-fetcher
```

2. **Levantar los Servicios con Docker Compose**

```bash
docker-compose up -d
```

Este comando inicia todos los servicios necesarios (backend, frontend, Redis, MariaDB) en contenedores Docker.

## Acceder al Proyecto

Una vez que los contenedores estén arriba y corriendo, podrás acceder a la aplicación web a través de:

- **Frontend**: http://localhost:8080/app
- **API Backend**: http://localhost:8080/api/v1

## Despliegue

Este proyecto está configurado para facilitar su despliegue con Docker, asegurando una instalación y ejecución consistentes en cualquier entorno.

## Contribuciones

Las contribuciones son bienvenidas. Si tienes alguna sugerencia para mejorar la aplicación, por favor, considera enviar un pull request o abrir un issue en el repositorio.

## Licencia

Este proyecto está licenciado bajo la MIT License - vea el archivo [LICENSE.md](LICENSE.md) para más detalles.

---

Desarrollado con ❤ por Hugo.
