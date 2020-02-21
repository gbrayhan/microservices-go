function busquedaMedicamentos() {
    let inputMedicamento = document.getElementById("busqueda-medicamento");
    let busqueda = inputMedicamento.value;
    let tablaMedicamentos = document.getElementById("tabla-medicamentos");

    busqueda = busqueda.trim();
    if (busqueda.length > 0) {
        $.ajax({
            url: '/medicamentos/show',
            type: 'POST',
            data: {
                cadena: busqueda
            },
            cache: false,
            beforeSend: function () {
                inputMedicamento.classList.add("is-loading")
            },
        }).done(function (response) {
            if (response.success) {
                inputMedicamento.classList.remove("is-loading");
                let tabla = "";
                response.medicamentos.forEach(function (medicamento) {
                    tabla += `
					<tr id="${medicamento.id}">
						<td>${medicamento.id}</td>
						<td id="codigo-${medicamento.id}">${medicamento.codigo_ean}</td>
						<td id="descripcion-${medicamento.id}">${medicamento.descripcion_articulo}</td>
						<td id="laboratorio-${medicamento.id}">${medicamento.laboratorio}</td>
						<td id="clave-sat-${medicamento.id}">${medicamento.clave_sat}</td>
						<td>
							<div class="flex-space-around">
								<i class="fas fa-pencil-alt hover-cursor" onclick="openModalMedicamento('edit', ${medicamento.id})"></i>
								<i class="fas fa-trash-alt hover-cursor" onclick="openModalMedicamentoRemove(${medicamento.id})"></i>
							</div>
						</td>
					</tr>
					`
                });

                tablaMedicamentos.innerHTML = tabla
            } else {
                tablaMedicamentos.innerHTML = ""
            }
        })
    } else {
        tablaMedicamentos.innerHTML = ""
    }
}

function saveMedicamento() {
    let codigoEan = document.getElementById("codigo-ean").value;
    let descripcion = document.getElementById("descripcion-articulo").value;
    let laboratorio = document.getElementById("laboratorio").value;
    let claveSat = document.getElementById("clave-sat").value;
    let notificacion = document.getElementById("notificacion");

    let json = {
        codigo_ean: codigoEan,
        descripcion_articulo: descripcion,
        laboratorio: laboratorio,
        clave_sat: claveSat
    };

    $.ajax({
        url: '/medicamentos/new',
        type: 'POST',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify(json),
        cache: false
    }).done(function (response) {
        if (response.success) {
            notificacion.classList.add("has-text-success");
            notificacion.innerText = response.message;
            setTimeout(function () {
                closeModalMedicamento('form');
                cleanModal();
                notificacion.classList.remove("has-text-success");
                notificacion.innerText = "";
            }, 3000)
        } else {
            notificacion.classList.add("has-text-danger");
            notificacion.innerText = response.message;
            setTimeout(function () {
                closeModalMedicamento('form');
                cleanModal();
                notificacion.classList.remove("has-text-danger");
                notificacion.innerText = "";
            }, 3000)
        }
    })
}

function uploadMedicamento(idStr) {
    let codigoEan = document.getElementById("codigo-ean").value;
    let descripcion = document.getElementById("descripcion-articulo").value;
    let laboratorio = document.getElementById("laboratorio").value;
    let claveSat = document.getElementById("clave-sat").value;
    let notificacion = document.getElementById("notificacion");

    let id = parseInt(idStr);

    let json = {
        id: id,
        codigo_ean: codigoEan,
        descripcion_articulo: descripcion,
        laboratorio: laboratorio,
        clave_sat: claveSat
    };

    $.ajax({
        url: '/medicamentos/update',
        type: 'POST',
        dataType: 'json',
        contentType: 'application/json',
        data: JSON.stringify(json),
        cache: false
    }).done(function (response) {
        if (response.success) {
            notificacion.classList.add("has-text-success");
            notificacion.innerText = response.message;
            setTimeout(function () {
                closeModalMedicamento('form');
                cleanModal();
                notificacion.classList.remove("has-text-success");
                notificacion.innerText = "";

                document.getElementById("codigo-" + id).innerText = codigoEan;
                document.getElementById("descripcion-" + id).innerText = descripcion;
                document.getElementById("laboratorio-" + id).innerText = laboratorio;
                document.getElementById("clave-sat-" + id).innerText = claveSat;
            }, 3000)
        } else {
            notificacion.classList.add("has-text-danger");
            notificacion.innerText = response.message;
            setTimeout(function () {
                closeModalMedicamento('form');
                cleanModal();
                notificacion.classList.remove("has-text-danger");
                notificacion.innerText = "";
            }, 3000)
        }
    })
}

function openModalMedicamento(origen, id) {
    if (origen === "new") {
        document.getElementById("modal-title").innerText = "Nuevo medicamento";
        let button = document.getElementById("modal-btn-success");
        button.innerText = "Agregar";
        button.setAttribute("onclick", "saveMedicamento()");

        cleanModal()
    } else if (origen === "edit") {
        document.getElementById("codigo-ean").value = document.getElementById("codigo-" + id).innerText;
        document.getElementById("descripcion-articulo").value = document.getElementById("descripcion-" + id).innerText;
        document.getElementById("laboratorio").value = document.getElementById("laboratorio-" + id).innerText;
        document.getElementById("clave-sat").value = document.getElementById("clave-sat-" + id).innerText;
        document.getElementById("modal-title").innerText = "Editar medicamento";
        let button = document.getElementById("modal-btn-success");
        button.innerText = "Actualizar";
        button.setAttribute("onclick", "uploadMedicamento('" + id + "')")
    }

    let modalMedicamentos = document.getElementById("modal-medicamentos-form");
    modalMedicamentos.classList.add("is-active")
}

function openModalMedicamentoRemove(id) {
    document.getElementById("modal-medicamentos-remove").classList.add("is-active");
    document.getElementById("btn-remove-medicamento").setAttribute("onclick", "removeMedicamento('" + id + "')")
}

function removeMedicamento(id) {
    let notificacion = document.getElementById("notificacion-remove");

    $.ajax({
        url: '/medicamentos/delete',
        type: 'POST',
        data: {
            id: id
        }
    }).done(function (response) {
        if (response.success) {
            notificacion.classList.add("has-text-success");
            notificacion.innerText = response.message;
            setTimeout(function () {
                closeModalMedicamento('remove');
                notificacion.classList.remove("has-text-success");
                notificacion.innerText = "";

                document.getElementById(id).remove();
            }, 3000)
        } else {
            notificacion.classList.add("has-text-danger");
            notificacion.innerText = response.message;
            setTimeout(function () {
                closeModalMedicamento('remove');
                notificacion.classList.remove("has-text-danger");
                notificacion.innerText = "";
            }, 3000)
        }
    })
}

function closeModalMedicamento(modal) {
    let modalMedicamentos = document.getElementById("modal-medicamentos-" + modal);
    modalMedicamentos.classList.remove("is-active")
}

function cleanModal() {
    document.getElementById("codigo-ean").value = "";
    document.getElementById("descripcion-articulo").value = "";
    document.getElementById("laboratorio").value = "";
    document.getElementById("clave-sat").value = "";
}