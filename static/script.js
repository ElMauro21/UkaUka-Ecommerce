document.addEventListener('DOMContentLoaded', () => {
    // ===== Form Toggles =====
    const loginForm = document.getElementById('login');
    const registerForm = document.getElementById('register');
    const recoverForm = document.getElementById('recover');
    const subMenu = document.getElementById('subMenu');
    const toggleButton = document.getElementById('toggleMenuBtn');
    const changeForm = document.getElementById('passwords');
    const profileForm = document.getElementById('profile');

    window.change = function () {
        profileForm?.style && (profileForm.style.display = 'none');
        changeForm?.style && (changeForm.style.display = 'block');
    };

    window.profile = function () {
        profileForm?.style && (profileForm.style.display = 'block');
        changeForm?.style && (changeForm.style.display = 'none');
    };

    window.login = function () {
        loginForm?.style && (loginForm.style.display = 'block');
        registerForm?.style && (registerForm.style.display = 'none');
        recoverForm?.style && (recoverForm.style.display = 'none');
    };

    window.register = function () {
        loginForm?.style && (loginForm.style.display = 'none');
        registerForm?.style && (registerForm.style.display = 'block');
        recoverForm?.style && (recoverForm.style.display = 'none');
    };

    window.recover = function () {
        loginForm?.style && (loginForm.style.display = 'none');
        recoverForm?.style && (recoverForm.style.display = 'block');
    };

    window.recoverLogin = function () {
        loginForm?.style && (loginForm.style.display = 'block');
        recoverForm?.style && (recoverForm.style.display = 'none');
    };

    if (toggleButton && subMenu) {
        toggleButton.addEventListener('click', event => {
            event.stopPropagation();
            subMenu.classList.toggle('open-menu');
        });

        document.addEventListener('click', event => {
            if (
                !subMenu.contains(event.target) &&
                !toggleButton.contains(event.target)
            ) {
                subMenu.classList.remove('open-menu');
            }
        });
    }

    setupAdminProductForm();
});

// ========== HTMX afterSwap Listener ==========
document.body.addEventListener('htmx:afterSwap', event => {
    if (event.detail.target.id === 'flash') {
        const flash = document.getElementById('flash');
        const type = flash?.getAttribute('data-type');

        if (type === 'success') {
            const form = document.getElementById('admin-products');
            form?.reset();

            const addButton = document.getElementById('add-button');
            if (addButton) addButton.style.display = 'inline-block';

            const hiddenId = form?.querySelector('[name="product-id"]');
            if (hiddenId) hiddenId.remove();
        }
    }

    if (event.detail.target.id === 'admin-products') {
        setupAdminProductForm();
    }
});

// ========== Product Form Handling ==========
function setupAdminProductForm() {
    const form = document.getElementById('admin-products');
    const addButton = document.getElementById('add-button');

    if (!form) return;

    window.fillProductForm = select => {
        const option = select.options[select.selectedIndex];
        if (!option.value) {
            form.reset();
            if (addButton) addButton.style.display = 'inline-block';

            const hiddenId = form.querySelector('[name="product-id"]');
            if (hiddenId) hiddenId.remove();
            return;
        }

        form.querySelector('[name="product-name"]').value = option.dataset.name;
        form.querySelector('[name="product-description"]').value =
            option.dataset.description;
        form.querySelector('[name="product-weight"]').value =
            option.dataset.weight;
        form.querySelector('[name="product-size"]').value = option.dataset.size;
        form.querySelector('[name="product-price"]').value =
            option.dataset.price;
        form.querySelector('[name="product-quantity"]').value =
            option.dataset.quantity;
        form.querySelector('[name="product-image"]').value =
            option.dataset.image;
        form.querySelector('[name="product-image-two"]').value =
            option.dataset.image2;

        let hiddenInput = form.querySelector('[name="product-id"]');
        if (!hiddenInput) {
            hiddenInput = document.createElement('input');
            hiddenInput.type = 'hidden';
            hiddenInput.name = 'product-id';
            form.appendChild(hiddenInput);
        }

        hiddenInput.value = option.value;

        if (addButton) addButton.style.display = 'none';
    };
}

// Price formatted
document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.price').forEach(el => {
        el.textContent = formatCOP(el.dataset.price);
    });
});

// Modal product

document.addEventListener('DOMContentLoaded', () => {
    const modalWrapper = document.getElementById('modal-wrapper');
    const modalOverlay = document.querySelector('.modal-overlay');
    const closeButton = document.getElementById('close');
    const openButtons = document.querySelectorAll('.open-modal');

    openButtons.forEach(button => {
        button.addEventListener('click', () => {
            const productID = button.dataset.id;
            const name = button.dataset.name;
            const description = button.dataset.description;
            const quantity = parseInt(button.dataset.quantity);
            const price = formatCOP(button.dataset.price);
            const image = button.dataset.image;

            document.querySelector(
                '#modal-container input[name="product-id"]'
            ).value = productID;
            document.querySelector(
                '#modal-container .modal-name h1'
            ).textContent = name;
            document.querySelector(
                '#modal-container .modal-description p'
            ).textContent = description;
            document.querySelector(
                '#modal-container .modal-price p'
            ).textContent = price;
            document.querySelector('#modal-container .modal-image img').src =
                image;
            stockAvailable.textContent = quantity;

            currentQty = 1;
            qtyDisplay.textContent = currentQty;
            quantityHidden.value = currentQty;

            modalWrapper.style.display = 'flex';
        });
    });

    if (modalOverlay) {
        modalOverlay.addEventListener('click', event => {
            if (event.target === modalOverlay) {
                modalWrapper.style.display = 'none';
            }
        });
    }

    if (closeButton) {
        closeButton.addEventListener('click', () => {
            modalWrapper.style.display = 'none';
        });
    }
});

let currentQty = 1;
const qtyDisplay = document.getElementById('qty-display');
const stockAvailable = document.getElementById('stock-available');
const quantityHidden = document.getElementById('quantity-hidden');

document.getElementById('qty-increase').addEventListener('click', () => {
    const max = parseInt(stockAvailable.textContent);
    if (currentQty < max) {
        currentQty += 1;
        qtyDisplay.textContent = currentQty;
        quantityHidden.value = currentQty;
    }
});

document.getElementById('qty-decrease').addEventListener('click', () => {
    if (currentQty > 1) {
        currentQty -= 1;
        qtyDisplay.textContent = currentQty;
        quantityHidden.value = currentQty;
    }
});
