document.addEventListener('DOMContentLoaded', function () {
    const orderIdInput = document.getElementById('orderIdInput');
    const fetchOrderButton = document.getElementById('fetchOrderButton');
    const createOrderButton = document.getElementById('createOrderButton');
    const deleteOrderButton = document.getElementById('deleteOrderButton');

    fetchOrderButton.addEventListener('click', fetchOrderDetails);
    createOrderButton.addEventListener('click', createOrder);
    deleteOrderButton.addEventListener('click', deleteOrder);
    // document.addEventListener('DOMContentLoaded', fetchOrderIds);

    
    fetch('/orders/all')
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch IDs');
            }
            return response.json();
        })
        .then(data => {
            console.log('IDs fetched successfully:', data);
            const savedOrdersList = document.getElementById('savedOrdersList');
            for (const id of data) {
                const orderIdItem = document.createElement('li');
                orderIdItem.textContent = "Order ID: " + id;
                savedOrdersList.appendChild(orderIdItem);
            }
        })
        .catch(error => {
            console.error('Error fetching IDs:', error);
        });

    function fetchOrderDetails() {
        const id = orderIdInput.value;
        if (id === '') {
            const errorDiv = document.getElementById('errorDiv');
            errorDiv.textContent = 'Input field not must be empty';
            return;
        };
        fetch(`/orders/id/${id}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to fetch order details');
                }
                return response.json();
            })
            .then(order => {
                const orderDetailsDiv = document.getElementById('orderDetails');
                orderDetailsDiv.innerHTML = '<pre>' + JSON.stringify(order, null, 2) + '</pre>';
            })
            .catch(error => {
                console.error('Error fetching order details:', error);
                const orderDetailsDiv = document.getElementById('orderDetails');
                orderDetailsDiv.innerHTML = '<p>Error fetching order details: ' + error.message + '</p>';
            });
    }

    function createOrder() {
        fetch(`/orders/new`, {
            method: 'POST',
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to create order');
            }
            return response.json();
        })
        .then(data => {
            console.log('Order created successfully:', data);
            const savedOrdersList = document.getElementById('savedOrdersList');
            const orderIdItem = document.createElement('li');
            orderIdItem.textContent = "Order ID: " + data['order_uid'];
            savedOrdersList.appendChild(orderIdItem);
        })
        .catch(error => {
            console.error('Error creating order:', error);
            const errorDiv = document.getElementById('errorDiv');
            errorDiv.textContent = 'Error creating order: ' + error.message;
        });
    }      

    function deleteOrder() {
        const id = orderIdInput.value;
        if (id === '') {
            const errorDiv = document.getElementById('errorDiv');
            errorDiv.textContent = 'Input field not must be empty';
            return;
        };
        fetch(`/orders/id/${id}`, {
            method: 'DELETE',
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to delete order');
            }
        })
        .then(data => {
            const savedOrdersList = document.getElementById('savedOrdersList');
            const orderItems = savedOrdersList.getElementsByTagName('li');
            for (let i = 0; i < orderItems.length; i++) {
                if (orderItems[i].textContent.includes(`Order ID: ${id}`)) {
                    savedOrdersList.removeChild(orderItems[i]);
                    break;
                }
            }
            const orderDetailsDiv = document.getElementById('orderDetails');
            orderDetailsDiv.innerHTML = '<p>Order deleted successfully</p>';
        })
        .catch(error => {
            console.error('Error deleting order:', error);
            const orderDetailsDiv = document.getElementById('orderDetails');
            orderDetailsDiv.innerHTML = '<p>Error deleting order: ' + error.message + '</p>';
        });
    }    
})