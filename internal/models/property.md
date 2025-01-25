## Visión General

El struct `Property` define la estructura de datos para representar propiedades inmobiliarias obtenidas mediante scraping. Actúa como contenedor para toda la información relevante de cada propiedad.

---

## Estructura del Modelo

### Definición completa

```go
type Property struct {
    Title     string
    Price     string
    Location  string
    M2        string
    Bedrooms  string
    Bathrooms string
    Link      string
}
```

---

## Campos Detallados

| Campo         | Tipo     | Descripción                           | Ejemplo                                              |
| ------------- | -------- | ------------------------------------- | ---------------------------------------------------- |
| **Title**     | `string` | Nombre/título de la propiedad         | "Casa 3 dormitorios condominio cerrado"              |
| **Price**     | `string` | Precio en UF **con formato original** | "UF 4.950"                                           |
| **Location**  | `string` | Ubicación específica                  | "Las Condes, Sector San Carlos de Apoquindo"         |
| **M2**        | `string` | Metros cuadrados **con unidad**       | "120 m²"                                             |
| **Bedrooms**  | `string` | Cantidad de dormitorios               | "3 dormitorios"                                      |
| **Bathrooms** | `string` | Cantidad de baños                     | "2 baños"                                            |
| **Link**      | `string` | URL única de la propiedad             | "<https://www.portalinmobiliario.com/MLC-18293723>..." |

---

## Propósito y Diseño

1. **Consistencia con fuente original**: Mantiene formato crudo del sitio web para:

   - Facilitar debugging
   - Permitir procesamiento posterior
   - Mantener fidelidad de datos originales

2. **Flexibilidad en tipos de datos**:

   - Todos los campos como `string` para manejar:
     - Valores numéricos con unidades ("120 m²")
     - Formatos variables ("UF3.500" vs "3.500 UF")
     - Datos incompletos/mal formados

3. **Relación con scraping**:
   - Estructura refleja directamente los datos extraídos
   - Campos mapean 1:1 con elementos HTML scrapeados

---

## Ejemplo de Uso

```go
prop := models.Property{
    Title:     "Casa moderna 3D/2B con piscina",
    Price:     "UF 5.200",
    Location:  "Lo Barnechea, Av. El Rodeo",
    M2:        "185 m²",
    Bedrooms:  "3 dormitorios",
    Bathrooms: "2 baños",
    Link:      "https://www.portalinmobiliario.com/MLC-123456",
}
```

---

## Relación con Otros Componentes

1. **Scraper**:

   ```mermaid
   graph LR
   A[HTML] --> B[Scraper] --> C[Property]
   ```

   - El scraper llena cada campo extrayendo datos específicos del HTML

2. **Análisis**:
   - Campos `Price` y `M2` son procesados por:

   ```go
   func extraerValorUF(precio string) int
   func extraerValorM2(m2 string) int
   ```

   - Transforman strings a valores numéricos para cálculos
