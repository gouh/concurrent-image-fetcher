# Proyecto de Descarga Concurrente de Imágenes

Este proyecto es una aplicación web diseñada para descargar imágenes desde múltiples URLs de manera concurrente, permitiendo al usuario especificar el número deseado de descargas concurrentes. La aplicación backend está desarrollada en Go, aprovechando sus capacidades para operaciones concurrentes, mientras que el frontend se implementa en React, ofreciendo una interfaz de usuario interactiva y atractiva.

## Características

- Validación y manejo de URLs de entrada.
- Limitación configurable del número de descargas concurrentes.
- Sistema robusto para manejar errores durante las descargas.
- Indicadores de progreso de descarga en tiempo real.
- Interfaz de usuario React para interacción con el servicio.
- Contenerización con Docker para fácil despliegue y distribución.

## Requisitos Previos

Para construir y ejecutar esta aplicación, necesitarás:

- Go (versión 1.16 o superior)
- Node.js y npm (para el frontend React)
- Docker

## Instalación

### Backend (Go)

1. Clona el repositorio y navega al directorio del backend.
2. Ejecuta `go build` para compilar el proyecto.
3. Inicia el servidor ejecutando el binario compilado, por ejemplo, `./miAplicacion`.

### Frontend (React)

1. Navega al directorio del frontend.
2. Ejecuta `npm install` para instalar las dependencias.
3. Inicia la aplicación React con `npm start`.

### Uso de Docker

1. En la raíz del proyecto, construye la imagen de Docker con `docker build -t nombre-de-tu-imagen .`.
2. Ejecuta el contenedor usando `docker run -p 4000:4000 nombre-de-tu-imagen`, ajustando los puertos según sea necesario.

## Uso

### Cargar URLs de Imágenes

- Accede a la interfaz de usuario de React en `http://localhost:3000` (ajusta el puerto según tu configuración).
- Usa el formulario proporcionado para introducir o cargar las URLs de las imágenes que deseas descargar.

### Configuración de Descargas Concurrentes

- Utiliza la interfaz de configuración en la página para ajustar el número de descargas concurrentes permitidas.

## Contribuir

Si deseas contribuir a este proyecto, por favor considera:

- Crear un fork del repositorio.
- Crear una rama para tu característica (`git checkout -b feature/fooBar`).
- Hacer commit de tus cambios (`git commit -am 'Add some fooBar'`).
- Hacer push a la rama (`git push origin feature/fooBar`).
- Abrir una nueva Pull Request.

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - vea el archivo `LICENSE.md` para detalles.
