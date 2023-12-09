const getBtn = document.querySelector("#idBtn")
const getInput = document.querySelector("input")


const getData = () => {
    let id = getInput.value
    if (id <= 0) {
        console.log("ERROR!: Empty input...");
    } else {
        fetch(`http://localhost:8080/orders/${id}`).then(response => {
            return response.json()
        }).then((data) => {
            console.log(data)
            let orderData = data
            let cardItem = ''
            let card = document.getElementById('card')
            cardItem +=
                `
    <div style="border:2px solid #ccc;width: 300px; margin: 12px 50px; padding: 10px 20px">
    <h2>ORDER: ${orderData.id}</h2>
    <p>order_uid: ${orderData.order_uid}</p>
    <p>track_number: ${orderData.track_number}</p>
    <p>entry: ${orderData.entry}</p>
    
        <div style="border:1px solid #ccc;width: 200px;padding: 10px 20px">
        <h3>Delivery: </h3>
        <p>name: ${orderData.delivery.name}</p>
        <p>phone: ${orderData.delivery.phone}</p>
        <p>zip: ${orderData.delivery.zip}</p>
        <p>city: ${orderData.delivery.city}</p>
        <p>address: ${orderData.delivery.address}</p>
        <p>region: ${orderData.delivery.region}</p>
        <p>email: ${orderData.delivery.email}</p>
        </div>
        
        <div style="border:1px solid #ccc;width: 200px;padding: 10px 20px">
        <h3>Payment: </h3>
        <p>transaction: ${orderData.payment.transaction}</p>
        <p>request_id: ${orderData.payment.request_id}</p>
        <p>currency: ${orderData.payment.currency}</p>
        <p>provider: ${orderData.payment.provider}</p>
        <p>amount: ${orderData.payment.amount}</p>
        <p>payment_dt: ${orderData.payment.payment_dt}</p>
        <p>bank: ${orderData.payment.bank}</p>
        <p>delivery_cost: ${orderData.payment.delivery_cost}</p>
        <p>goods_total: ${orderData.payment.goods_total}</p>
        <p>custom_fee: ${orderData.payment.custom_fee}</p>
        </div>
        
        <div id="itemId" style="border:1px solid #ccc;width: 200px;padding: 10px 20px"></div>
        
    <p>locale: ${orderData.locale}</p>
    <p>internal_signature: ${orderData.internal_signature}</p>
    <p>customer_id: ${orderData.customer_id}</p>
    <p>delivery_service: ${orderData.delivery_service}</p>
    <p>shardkey: ${orderData.shardkey}</p>
    <p>sm_id: ${orderData.sm_id}</p>
    <p>date_created: ${orderData.date_created}</p>
    <p>oof_shard: ${orderData.oof_shard}</p>
    <p>created_at: ${orderData.created_at}</p>
    </div>
    `
            card.insertAdjacentHTML('afterbegin', cardItem);

            let item = ''
            let itemId = document.getElementById('itemId')

            orderData.items.forEach(
                (itemUnit) => {
                    item +=
                        `
               <h3>Item: </h3>
        <p>chrt_id: ${itemUnit.chrt_id}</p>
        <p>track_number: ${itemUnit.track_number}</p>
        <p>price: ${itemUnit.price}</p>
        <p>rid: ${itemUnit.rid}</p>
        <p>name: ${itemUnit.name}</p>
        <p>sale: ${itemUnit.sale}</p>
        <p>size: ${itemUnit.size}</p>
        <p>total_price: ${itemUnit.total_price}</p>
        <p>nmId: ${itemUnit.nmId}</p>
        <p>brand: ${itemUnit.brand}</p>
        <p>status: ${itemUnit.status}</p>
        <p>updated_at: ${itemUnit.updated_at}</p>
        <p>created_at: ${itemUnit.created_at}</p>
                 `
                }
            )
            itemId.insertAdjacentHTML("afterbegin", item)


        })
    }
};

getBtn.addEventListener('click', getData)
