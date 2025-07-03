FROM node:20

WORKDIR /app

# Copia dependências e instala
COPY package*.json ./
RUN npm install

# Copia código-fonte
COPY . .

# Mantém o container "aberto" para desenvolvimento
CMD ["tail", "-f", "/dev/null"]
