# Nexus Shell (nxsh)

**Una shell de línea de comandos moderna y consciente de los datos, escrita en Go.**

Nexus Shell (nxsh) revoluciona la forma en que interactúas con los datos en la terminal. A diferencia de las shells tradicionales que tratan todo como texto plano, `nxsh` detecta y parsea automáticamente las salidas JSON, convirtiéndolas en objetos de primera clase que puedes consultar, filtrar y transformar de forma nativa.

## Características Principales

*   **Detección Automática de JSON:** `nxsh` inspecciona la salida de los comandos. Si es JSON válido, lo trata como datos estructurados, no como una simple cadena de texto.
*   **Manipulación de Datos Nativa:** Usa los comandos internos `get`, `where`, y `select` para consultar, filtrar y transformar datos JSON de forma intuitiva.
*   **Pipelines Potentes:** Encadena comandos como en cualquier shell, pero con la capacidad de pasar objetos de datos estructurados entre comandos internos, no solo texto.
*   **Variables y Estado:** Usa `let` para guardar la salida de cualquier comando en una variable y reutilizarla más tarde.
*   **REPL Interactivo:** Una experiencia de terminal moderna con historial de comandos persistente y un prompt dinámico y con colores.

## Instalación

Asegúrate de tener Go instalado (versión 1.18 o superior). Puedes instalar `nxsh` directamente con:

```bash
go install github.com/soyunomas/nxsh@latest
```

***Nota:*** *El comando anterior compilará e instalará el binario `nxsh` en el directorio de binarios de Go (normalmente en `$HOME/go/bin`). Puedes ejecutar el programa directamente usando su ruta completa:*

```bash
~/go/bin/nxsh
```

*Para poder invocarlo de forma más cómoda (escribiendo solo `nxsh`), puedes añadir su directorio a la variable de entorno `$PATH` de tu sistema. Para hacerlo de forma permanente, ejecuta:*

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc && source ~/.bashrc
```
*(Si usas Zsh, reemplaza `~/.bashrc` por `~/.zshrc`)*.

O clona el repositorio y compila localmente:

```bash
git clone https://github.com/soyunomas/nxsh.git
cd nxsh
go build
./nxsh
```

## Ejemplos de Uso Detallados

Imaginemos que tenemos un archivo `users.json` con el siguiente contenido, y lo cargamos en una variable:

```json
[
  {"id": 1, "name": "Alice", "role": "admin", "isActive": true, "location": {"city": "New York", "country": "USA"}, "age": 28},
  {"id": 2, "name": "Bob", "role": "developer", "isActive": false, "location": {"city": "London", "country": "UK"}, "age": 35},
  {"id": 3, "name": "Charlie", "role": "developer", "isActive": true, "location": {"city": "New York", "country": "USA"}, "age": 42},
  {"id": 4, "name": "Diana", "role": "guest", "isActive": true, "location": {"city": "Tokyo", "country": "Japan"}, "age": 25}
]
```

```shell
nxsh > let users = cat users.json
```

Ahora podemos manipular la variable `users` con los comandos internos.

### `get`: Extraer campos

El comando `get` te permite extraer valores de objetos o arrays de objetos.

**Ejemplo 1: Obtener los nombres de todos los usuarios**

```shell
nxsh > users | get .name
```**Salida:**
```json
[
  "Alice",
  "Bob",
  "Charlie",
  "Diana"
]
```

**Ejemplo 2: Obtener las ciudades de todos los usuarios (campo anidado)**
```shell
nxsh > users | get .location.city
```
**Salida:**
```json
[
  "New York",
  "London",
  "New York",
  "Tokyo"
]
```

### `where`: Filtrar arrays

El comando `where` te permite filtrar un array de objetos basándote en una condición.

**Ejemplo 1: Encontrar todos los usuarios que son "developer"**

```shell
nxsh > users | where .role == "developer"
```
**Salida:**
```json
[
  { "id": 2, "name": "Bob", ... },
  { "id": 3, "name": "Charlie", ... }
]
```

**Ejemplo 2: Encontrar usuarios activos mayores de 30 años**
```shell
nxsh > users | where .isActive == true | where .age > 30
```
**Salida:**
```json
[
  { "id": 3, "name": "Charlie", "role": "developer", "isActive": true, "location": {...}, "age": 42 }
]
```

### `select`: Remodelar objetos

El comando `select` te permite crear nuevos objetos, quedándote solo con los campos que te interesan.

**Ejemplo 1: Obtener solo el nombre y el rol de cada usuario**
```shell
nxsh > users | select .name .role
```
**Salida:**```json
[
  { "name": "Alice", "role": "admin" },
  { "name": "Bob", "role": "developer" },
  { "name": "Charlie", "role": "developer" },
  { "name": "Diana", "role": "guest" }
]
```

**Ejemplo 2: Combinar todo para obtener el nombre y la ciudad de los usuarios de USA**
```shell
nxsh > users | where .location.country == "USA" | select .name .location.city
```
**Salida:**
```json
[
  { "name": "Alice", "city": "New York" },
  { "name": "Charlie", "city": "New York" }
]
```

## Caso de Uso Real: API de GitHub

Encuentra los nombres y el número de estrellas de los repositorios de Google que no son forks.

```shell
# 1. Llama a la API de GitHub y guarda la respuesta JSON en una variable.
nxsh > let repos = curl -s "https://api.github.com/users/google/repos?per_page=10&sort=pushed"

# 2. Usa 'where' para filtrar la lista, quedándote solo con los que NO son forks.
nxsh > let originales = repos | where .fork == false

# 3. De esos, usa 'select' para quedarte solo con el nombre y las estrellas.
nxsh > originales | select .name .stargazers_count

# Resultado: Un nuevo array JSON con la información precisa que querías.
[
  {
    "name": "gvisor",
    "stargazers_count": 23456
  },
  {
    "name": "go-github",
    "stargazers_count": 9876
  },
  ...
]
```

## Hoja de Ruta del Proyecto

### ✅ Hitos Alcanzados

-   `[x]` **Arquitectura del Evaluador:** Refactorización completa a un evaluador que recorre el AST directamente.
-   `[x]` **REPL Moderno:** Shell interactiva con historial y prompt dinámico.
-   `[x]` **Conciencia de Datos:** Detección y parseo automático de JSON.
-   `[x]` **Variables:** Implementación de `let` y un entorno de variables persistente.
-   `[x]` **Pipelines Inteligentes:** Implementación de `|` que maneja tanto texto como objetos.
-   `[x]` **Tríada de Datos Completa:** Implementación de los comandos `get`, `where` y `select`.

### 📝 Hoja de Ruta (TODO)

-   `[ ]` **Mejoras del Lenguaje:**
    -   `[ ]` Soportar la sintaxis de expansión de variables `$variable`.
    -   `[ ]` Añadir soporte para tipos de datos numéricos (int, float) y operaciones aritméticas.
-   `[ ]` **Control de Flujo:**
    -   `[ ]` Implementar condicionales `if/else`.
    -   `[ ]` Implementar bucles `for` para iterar sobre arrays.
-   `[ ]` **Funciones:**
    -   `[ ]` Permitir funciones definidas por el usuario (`def`).
-   `[ ]` **Ecosistema y Calidad de Vida:**
    -   `[ ]` Añadir un modo no interactivo (estilo `jq`) con un flag `-c`.
    -   `[ ]` Implementar manejo de errores avanzado con `try/catch`.
    -   `[ ]` Añadir soporte para un archivo de configuración (`~/.nxshrc`).
-   `[ ]` **Librería Estándar:**
    -   `[ ]` Expandir el conjunto de comandos internos para tareas comunes (archivos, red, etc.).

## Licencia

Este proyecto está bajo la Licencia MIT.
