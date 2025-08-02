document.addEventListener('DOMContentLoaded', () => {
    // ===== Form Toggles =====
    const loginForm = document.getElementById('login');
    const registerForm = document.getElementById('register');
    const recoverForm = document.getElementById('recover');
    const subMenu = document.getElementById('subMenu');
    const toggleButton = document.getElementById('toggleMenuBtn');
    const changeForm = document.getElementById('passwords');
    const profileForm = document.getElementById('profile');
    const deleteUserForm = document.getElementById('eliminate');

    window.change = function () {
        profileForm?.style && (profileForm.style.display = 'none');
        changeForm?.style && (changeForm.style.display = 'block');
        deleteUserForm?.style && (deleteUserForm.style.display = 'none');
    };

    window.profile = function () {
        profileForm?.style && (profileForm.style.display = 'block');
        changeForm?.style && (changeForm.style.display = 'none');
        deleteUserForm?.style && (deleteUserForm.style.display = 'none');
    };

    window.eliminate = function () {
        profileForm?.style && (profileForm.style.display = 'none');
        changeForm?.style && (changeForm.style.display = 'none');
        deleteUserForm?.style && (deleteUserForm.style.display = 'block');
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
        el.textContent = formatCOP(el.textContent);
    });
});

function applyPriceFormatting() {
    document.querySelectorAll('.price').forEach(el => {
        const raw = el.dataset.raw;
        el.textContent = formatCOP(raw);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    applyPriceFormatting();
});

document.body.addEventListener('htmx:afterSwap', () => {
    applyPriceFormatting();
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

    let currentQty = 1;
    const qtyDisplay = document.getElementById('qty-display');
    const stockAvailable = document.getElementById('stock-available');
    const quantityHidden = document.getElementById('quantity-hidden');

    const increaseBtn = document.getElementById('qty-increase');
    const decreaseBtn = document.getElementById('qty-decrease');

    if (
        increaseBtn &&
        decreaseBtn &&
        qtyDisplay &&
        stockAvailable &&
        quantityHidden
    ) {
        increaseBtn.addEventListener('click', () => {
            const max = parseInt(stockAvailable.textContent);
            if (currentQty < max) {
                currentQty += 1;
                qtyDisplay.textContent = currentQty;
                quantityHidden.value = currentQty;
            }
        });

        decreaseBtn.addEventListener('click', () => {
            if (currentQty > 1) {
                currentQty -= 1;
                qtyDisplay.textContent = currentQty;
                quantityHidden.value = currentQty;
            }
        });
    }
});

// Safari bug screen solution
function setViewportHeight() {
    let vh = window.innerHeight * 0.01;
    document.documentElement.style.setProperty('--vh', `${vh}px`);
}

setViewportHeight();

let resizeTimer;
window.addEventListener('resize', () => {
    clearTimeout(resizeTimer);
    resizeTimer = setTimeout(setViewportHeight, 0);
});

// Menus

document.addEventListener('DOMContentLoaded', () => {
    const toggleButton = document.getElementById('toggleMenuBtn');
    const subMenu = document.getElementById('subMenu');
    const dropdownToggle = document.getElementById('menu-toggle');
    const dropdownContent = document.getElementById('dropdown-content');
    const dropdown = document.querySelector('.dropdown');

    // ----- Submenu toggle -----
    if (toggleButton && subMenu) {
        toggleButton.addEventListener('click', event => {
            event.stopPropagation();
            subMenu.classList.toggle('open-menu');
            dropdownContent?.classList.remove('open-dropdown');
        });
    }

    // ----- Dropdown toggle -----
    if (dropdownToggle && dropdownContent) {
        dropdownToggle.addEventListener('click', event => {
            event.stopPropagation();
            dropdownContent.classList.toggle('open-dropdown');
            subMenu?.classList.remove('open-menu');
        });
    }

    document.addEventListener('click', event => {
        if (!dropdown?.contains(event.target)) {
            dropdownContent?.classList.remove('open-dropdown');
        }

        if (
            !subMenu?.contains(event.target) &&
            !toggleButton?.contains(event.target)
        ) {
            subMenu?.classList.remove('open-menu');
        }
    });
});

// Terms and conditions

document.addEventListener('DOMContentLoaded', function () {
    const modalWrapper = document.getElementById('modal-wrapper-terms');
    const modalOverlay = document.querySelector('.modal-overlay-terms');
    const openModalLink = document.getElementById('open-modal-link');
    const closeModalBtn = document.getElementById('close-modal-terms'); // or close-modal-terms

    if (openModalLink && modalWrapper) {
        openModalLink.addEventListener('click', function (event) {
            event.preventDefault();
            modalWrapper.style.display = 'flex';
        });
    }

    if (closeModalBtn && modalWrapper) {
        closeModalBtn.addEventListener('click', function () {
            modalWrapper.style.display = 'none';
        });
    }

    if (modalOverlay) {
        modalOverlay.addEventListener('click', function (event) {
            if (event.target === modalOverlay) {
                modalWrapper.style.display = 'none';
            }
        });
    }
});

function soundActivator() {
    const video = document.getElementById('videoStopMotion');
    video.muted = false;
    video.play();
}

function fillTransactionForm(select) {
    const option = select.options[select.selectedIndex];

    if (!option.value) {
        document.getElementById('transaction-form').reset();
        return;
    }

    document.getElementById('reference-code').value =
        option.dataset.referenceCode;
    document.getElementById('total-amount').value = option.dataset.totalAmount;
    document.getElementById('created-at').value = option.dataset.createdAt;
    document.getElementById('user-name').value = option.dataset.userName;
    document.getElementById('user-surname').value = option.dataset.userSurname;

    document.getElementById('shipping-name').value =
        option.dataset.shippingName;
    document.getElementById('shipping-id').value = option.dataset.shippingId;
    document.getElementById('shipping-phone').value =
        option.dataset.shippingPhone;
    document.getElementById('shipping-email').value =
        option.dataset.shippingEmail;
    document.getElementById('shipping-address').value =
        option.dataset.shippingAddress;

    const productsStr = option.dataset.products;

    const products = productsStr.split(',').map(product => {
        const parts = product.trim().split(' ');

        const name = parts.slice(0, -2).join(' ');
        const quantity = parseInt(
            parts[parts.length - 2].replace('(', '').replace(')', '')
        );
        const price = parseFloat(parts[parts.length - 1].replace('$', ''));

        return {
            Name: name,
            Quantity: quantity,
            Price: price,
        };
    });

    const productsList = document.getElementById('products-list');
    productsList.innerHTML = '';

    products.forEach(product => {
        const li = document.createElement('li');
        li.textContent = `${product.Name} | Cantidad: ${product.Quantity} | Precio: $${product.Price}`;
        productsList.appendChild(li);
    });
}
