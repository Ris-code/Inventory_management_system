console.log('{{.}}');

const backendData = JSON.parse('{{.}}');

// Access the properties as needed
const clubInfo = backendData.Club;
const items = backendData.Items;
const quantities = backendData.Quantity;

console.log(clubInfo);
console.log(items);
console.log(quantities);

// Function to create a product card
function createProductCard(item, quantity) {
    const cardContainer = document.getElementById("product-list-container");

    const card = document.createElement("div");
    card.className = "col-md-4";
    card.innerHTML = `
        <section class="panel">
            <div class="pro-img-box">
                <span>
                    <h3 class="pro-title">${item}</h3>
                    <p class="pro-quantity">Quantity: ${quantity}</p>
                </span>
                <a href="#" class="adtocart">
                    <i class="fa fa-shopping-cart"></i>
                </a>
            </div>
        </section>
    `;

    cardContainer.appendChild(card);
}

// Create product cards using data from the backend
for (let i = 0; i < items.length; i++) {
    createProductCard(items[i], quantities[i]);
}
