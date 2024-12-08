# Usa una imagen base de Ubuntu
FROM metasploitframework/metasploit-framework:latest

RUN apk update && apk add --no-cache \
    tcpdump \
    tshark

# Copia el código fuente de Go al contenedor
WORKDIR /app
COPY . .

# Crea la carpeta donde la aplicación guardará los resultados
RUN mkdir -p /app/results

# Define la carpeta como un volumen que se puede montar desde el host
VOLUME /app/results

# Compila el proyecto de Go
RUN go build -o server ./cmd

# Instala los requerimientos de Python
COPY scripts/ /app/
RUN pip3 install -r /app/requirements.txt

# Exponer el puerto en el que el servidor de Go correrá
EXPOSE 8080

# Comando para iniciar Metasploit y el servidor de Go
ENTRYPOINT ["/bin/sh", "-c", "ruby /usr/src/metasploit-framework/msfrpcd -U msf -P dL0rHLep -p 55552 -S false -f & ./server"]