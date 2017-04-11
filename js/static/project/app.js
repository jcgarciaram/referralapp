window.onload = function () {
    new Vue({
      el: '#cake-table',
      data: {
        cakes: [
          { id: '1', guestName: 'Pipo Antonio Garcia Marroquin', message: 'Happy Birthday!', flavor: 'Strawberry'},
          { id: '2', guestName: 'Rocco Marroquin', message: 'Feliz Cumpleaños!', flavor: 'Vanilla'}
        ]
      }
    });

    new Vue({
      el: '#app',
      data: {
        message: 'Hello Vue!'
      }
    });
}