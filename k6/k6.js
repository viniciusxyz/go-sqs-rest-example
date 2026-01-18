import http from 'k6/http';

export default function () {

    var requestBody = JSON.stringify({
        contaOrigem: "100000",
        contaDestino: "400001",
        valor: 1000
    });

    const params = {
        headers: {
        'Content-Type': 'application/json',
        },
    };

    var response = http.post("http://localhost:8080/deposito", requestBody, params);
    console.info(response.body)
}