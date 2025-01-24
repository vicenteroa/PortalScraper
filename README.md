# PortalScraper 🏠
![pORTALsCR](https://github.com/user-attachments/assets/f85b6729-4e4d-4ea2-abcc-0abeb61b7cf3)


PortalScraper es un scraper escrito en Go que extrae información de propiedades en venta desde el portal [Portal Inmobiliario](https://www.portalinmobiliario.com). Este proyecto está diseñado para obtener detalles como el título, precio, ubicación, metros cuadrados, número de dormitorios, baños y enlaces de las propiedades listadas.

---

## Características ✨

- **Extracción de datos**: Obtén detalles clave de propiedades en venta.
- **Paginación**: Soporte para scrapear múltiples páginas de resultados.
- **Formateo de datos**: Los datos se limpian y formatean para su fácil uso.
- **Configuración flexible**: Ajusta el número máximo de páginas a scrapear.
- **Respetuoso con el servidor**: Incluye un delay entre solicitudes para evitar sobrecargar el servidor.

---

## Requisitos 📋

- Go 1.20 o superior.
- Dependencias externas:
  - `github.com/PuerkitoBio/goquery` (para el análisis de HTML).

---

## Instalación 🛠️

1. Clona el repositorio:
   ```bash
   git clone https://github.com/vicenteroa/PortalScraper.git
   cd PortalScraper
   ```

2. Instala las dependencias:
   ```bash
   go mod tidy
   ```

3. Ejecuta el scraper:
   ```bash
   go run main.go
   ```

---

## Uso 🚀

El scraper está configurado para extraer datos de propiedades en venta en Las Condes, Santiago de Chile. Puedes modificar la URL base en el código para scrapear otras ubicaciones o tipos de propiedades.

### Ejemplo de salida:
```plaintext
Scrapeando página 1: https://www.portalinmobiliario.com/venta/casa/propiedades-usadas/las-condes-metropolitana?_PAGE=1

Total propiedades: 30

Propiedad #1:
Título: Casa en venta en Las Condes
Precio: $450,000,000
Ubicación: Las Condes, Santiago
m²: 250
Dormitorios: 4
Baños: 3
Enlace: https://www.portalinmobiliario.com/MLC-123456789
```

---

## Configuración ⚙️

- **URL base**: Modifica la variable `baseURL` en `main.go` para cambiar la ubicación o tipo de propiedad.
- **Número de páginas**: Ajusta la variable `maxPages` para scrapear más o menos páginas.
- **Delay entre solicitudes**: Cambia el valor de `time.Sleep(2 * time.Second)` para ajustar el tiempo de espera entre solicitudes.

---

## Estructura del código 🧩

- **`main.go`**: Contiene la lógica principal del scraper.
  - `scrapeMainPage`: Extrae datos de la página principal.
  - `cleanText`: Limpia y formatea los textos extraídos.
  - `extractLink`: Obtiene el enlace de la propiedad.
- **`Property` struct**: Almacena los detalles de cada propiedad.

---

## Contribución 🤝

¡Las contribuciones son bienvenidas! Si deseas mejorar el scraper, sigue estos pasos:

1. Haz un fork del repositorio.
2. Crea una rama con tu feature o fix: `git checkout -b mi-feature`.
3. Envía un pull request con tus cambios.

---

## Licencia 📜

Este proyecto está bajo la licencia **MIT**. Consulta el archivo [LICENSE](LICENSE) para más detalles.

---

## Advertencia ⚠️

Este scraper es solo para fines educativos. Asegúrate de cumplir con los términos de servicio de [Portal Inmobiliario](https://www.portalinmobiliario.com) y las leyes locales antes de usarlo en producción.

---

## Contacto 📧

Si tienes preguntas o sugerencias, no dudes en contactarme:  
📩 [tuemail@example.com](mailto:tuemail@example.com)  
🌐 [GitHub](https://github.com/tuusuario)

---

¡Gracias por usar **PortalScraper**! 🎉
