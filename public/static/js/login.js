$(document).ready(function () {

  $("#login").submit(function (e) {
    e.preventDefault();
    var username = $("#username").val();
    var password = $("#password").val();
    // Checking for blank fields.
    if (username == '' || password == '') {
      $('input[type="text"],input[type="username"]').css("border", "2px solid red");
      $('input[type="text"],input[type="password"]').css("box-shadow", "0 0 3px red");
    } else {
      var request = JSON.stringify({
        'username': username,
        'password': password
      });

      $.ajax({
        url: '/login',
        type: 'POST',
        contentType: "application/json",
        dataType: 'html',
        data: request,
        success: function (data) {
          var response = jQuery.parseJSON(data);

          console.log(data);
          console.log(response);

          if (response.success) {
            window.location.href = "/admin";
          }
          // $('#modal-contenido').text("Se ha guardado el calendario del ejercicio " + ejercicio + " en exito.");
          // $("#modal-accion").modal('show');
          // getCalendariosEmpresa(empresa);
        },
        error: function (msg) {
          console.log('error->', msg);
        }
      });

    }
  });

});