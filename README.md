# Correos Masivos

Este proyecto ofrece herramientas para el envío automatizado de correos electrónicos a gran escala. Facilita la comunicación con numerosos destinatarios de forma eficiente y controlada.

## Funcionalidades

- **Automatización:** Configura campañas y permite que el sistema realice los envíos automáticamente.
- **Personalización de Correos:** Adapta cada correo para incluir información específica del destinatario.

## Requisitos

Asegúrate de contar con lo siguiente en tu sistema:

- Go (Golang)
- Excelize library
- Base de datos SQL

## Configuración

Pasos para instalar y configurar el proyecto:

1. Clonar el repositorio localmente.
2. Ejecutar `go mod tidy` para instalar las dependencias.
3. Establecer las variables de entorno requeridas para la autenticación y envío de correos.

## Uso

Para enviar correos masivos, sigue estos comandos:

```shell
go build ./src/cmd
# Crear la tabla en la base de datos
./cmd.exe -create-table=true
# Importar contactos de un archivo Excel al formato 'NOMBRE CORREO' (solo dos columnas)
./cmd.exe -insert-excel=true -p "./path/excel.xlsx"
# Importar y enviar correos desde un archivo Excel al formato 'NOMBRE CORREO'
./cmd.exe -send-from-excel=true -p "./path/excel.xlsx"
# Enviar correos usando la configuración de variables de entorno
./cmd.exe
