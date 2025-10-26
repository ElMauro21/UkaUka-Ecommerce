function showFlashAlert() {
    const flash = document.getElementById('flash');
    if (flash) {
        const message = flash.dataset.message;
        const type = flash.dataset.type;

        if (message) {
            Swal.fire({
                toast: true,
                position: 'top-start',
                icon: type || 'info',
                title: message,
                showConfirmButton: false,
                timer: 3000,
                timerProgressBar: true,
            });
        }
    }
}

// Show flash on initial page load
document.addEventListener('DOMContentLoaded', showFlashAlert);

// Show flash after HTMX swaps in new content
document.body.addEventListener('htmx:afterSwap', function (e) {
    console.log('HTMX swapped something');
    showFlashAlert();
});
