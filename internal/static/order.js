document.addEventListener('DOMContentLoaded', function () {
    const fetchOrderButton = document.getElementById('fetchOrderButton');
    fetchOrderButton.addEventListener('click', fetchOrderDetails);

    function fetchOrderDetails() {
        const orderId = document.getElementById('orderIdInput').value;
        fetch(`/orders/id/${orderId}`)
            .then(response => response.json())
            .then(order => {
                const orderDetailsDiv = document.getElementById('orderDetails');
                orderDetailsDiv.innerHTML = ''; // Clear previous content
                if (order) {
                    orderDetailsDiv.innerHTML += `
                        <div class="section">
                            <h2>Order Details</h2>
                            <div class="details">
                                <p>Order ID: ${order.order_uid}</p>
                                <p>Track Number: ${order.track_number}</p>
                                <p>Entry: ${order.entry}</p>
                                <p>Locale: ${order.locale}</p>
                                <p>Internal Signature: ${order.internal_signature}</p>
                                <p>Customer ID: ${order.customer_id}</p>
                                <p>Delivery Service: ${order.delivery_service}</p>
                                <p>Shardkey: ${order.shardkey}</p>
                                <p>SmID: ${order.sm_id}</p>
                                <p>Date Created: ${order.date_created}</p>
                                <p>Oof Shard: ${order.oof_shard}</p>
                            </div>
                        </div>
                        <div class="section">
                            <h2>Delivery Details</h2>
                            <div class="details">
                                <p>Name: ${order.delivery.name}</p>
                                <p>Phone: ${order.delivery.phone}</p>
                                <p>Zip: ${order.delivery.zip}</p>
                                <p>City: ${order.delivery.city}</p>
                                <p>Address: ${order.delivery.address}</p>
                                <p>Region: ${order.delivery.region}</p>
                                <p>Email: ${order.delivery.email}</p>
                            </div>
                        </div>
                        <div class="section">
                            <h2>Items Details</h2>
                            ${order.items.map((item, index) => `
                                <div class="item-details">
                                    <h3>Item ${index + 1}</h3>
                                    <p>Chrt ID: ${item.chrt_id}</p>
                                    <p>Track Number: ${item.track_number}</p>
                                    <p>Price: ${item.price}</p>
                                    <p>Rid: ${item.rid}</p>
                                    <p>Name: ${item.name}</p>
                                    <p>Sale: ${item.sale}</p>
                                    <p>Size: ${item.size}</p>
                                    <p>Total Price: ${item.total_price}</p>
                                    <p>Nm ID: ${item.nm_id}</p>
                                    <p>Brand: ${item.brand}</p>
                                    <p>Status: ${item.status}</p>
                                </div>
                            `).join('')}
                        </div>
                        <div class="section">
                            <h2>Payment Details</h2>
                            <div class="details">
                                <p>Transaction: ${order.payment.transaction}</p>
                                <p>Request ID: ${order.payment.request_id}</p>
                                <p>Currency: ${order.payment.currency}</p>
                                <p>Provider: ${order.payment.provider}</p>
                                <p>Amount: ${order.payment.amount}</p>
                                <p>Payment Datetime: ${order.payment.payment_dt}</p>
                                <p>Bank: ${order.payment.bank}</p>
                                <p>Delivery Cost: ${order.payment.delivery_cost}</p>
                                <p>Goods Total: ${order.payment.goods_total}</p>
                                <p>Custom Fee: ${order.payment.custom_fee}</p>
                            </div>
                        </div>
                    `;
                } else {
                    orderDetailsDiv.innerHTML = '<p>Order not found</p>';
                }
            })
            .catch(error => {
                console.error('Error fetching order details:', error);
                const orderDetailsDiv = document.getElementById('orderDetails');
                orderDetailsDiv.innerHTML = '<p>Error fetching order details. Please try again later.</p>';
            });
    }
});
