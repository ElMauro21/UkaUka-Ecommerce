// Price formatted
function formatCOP(value) {
    const raw = parseFloat(value);
    return isNaN(raw)
        ? 'Precio no disponbible'
        : raw.toLocaleString('es-CO', {
              style: 'currency',
              currency: 'COP',
              minimumFractionDigits: 2,
          });
}
