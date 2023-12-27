# Correos Masivos

Este proyecto contiene un conjunto de herramientas y scripts que permiten el envío de correos electrónicos de forma masiva. Está diseñado para facilitar la comunicación con una gran cantidad de destinatarios de manera eficiente y automatizada.

## Características

- **Automatización del Envío:** Configura tus campañas de correo y deja que el sistema se encargue del resto.
- **Personalización:** Adapta los correos para cada destinatario incluyendo información relevante y personal.
- **Seguimiento:** Monitorea la entrega y la apertura de los correos para medir la efectividad de tus campañas.

## Requisitos Previos

Antes de comenzar, asegúrate de tener instalados los siguientes componentes en tu sistema:

- Go (Golang)
- Excelize library
- A SQL database

## Instalación

Para instalar y configurar el proyecto, sigue estos pasos:

1. Clona el repositorio en tu máquina local.
2. Instala las dependencias con `go mod tidy`.
3. Configura las variables de entorno necesarias para la autenticación y el envío de correos.

## Uso

Para comenzar a enviar correos masivos, ejecuta el script principal con los parámetros adecuados:

```shell
go build ./src/cmd
# para crear la tabla
./cmd.exe -create-table=true
# inserta los nombre y correos desde un archivo de excel a la base de datos. el excel debe tener el siguiente formato: NOMBRE	CORREO. y solo esas dos columnas
./cmd.exe -insert-excel=true -p "./path/excel.xlsx"
# envia los correos con la configuracion dada en las variables de entorno
./cmd.exe