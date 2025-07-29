FROM node:20-slim

RUN apt update && \
    apt install openssl procps -y && \
    npm install -g @nestjs/cli@11.0.6

WORKDIR /home/node/app

# Copia os arquivos de dependências primeiro
COPY package*.json ./

# Instala as dependências como root (padrão)
RUN npm install

# Copia o restante do código da aplicação
COPY . .
CMD [ "tail", "-f", "/dev/null" ]