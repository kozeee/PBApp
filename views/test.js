function test(x){
    let itemList =[{priceId:x,quantity:1}]
        Paddle.Checkout.Open(
        {settings:{displayMode:"overlay",},items:itemList})

}

async function checkemail(){
    let email = $('email').value
    let resp = await axios.get("/email/"+email)
    console.log(resp.data)
}