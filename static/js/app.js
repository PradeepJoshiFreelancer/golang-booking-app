
function Promt() {
    const toast = ({
        icon = "success",
        title = "",
        position = "top-end",
    }) => {
        const Toast = Swal.mixin({
            toast: true,
            position: position,
            icon: icon,
            title: title,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener("mouseenter", Swal.stopTimer);
                toast.addEventListener("mouseleave", Swal.resumeTimer);
            },
        });

        Toast.fire();
    };
    const success = ({ title = "", msg = "", footer = "" }) => {
        Swal.fire({
            icon: "success",
            title: title,
            text: msg,
            footer: footer,
        });
    };
    const error = ({ title = "", msg = "", footer = "" }) => {
        Swal.fire({
            icon: "error",
            title: title,
            text: msg,
            footer: footer
        });
    };
    const custom = async ({ icon = "", title = "", msg = "", showConfirmButton = true, callback, willOpen, preConfirm }) => {
        const { value: results } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            showConfirmButton: showConfirmButton,
            backdrop: true,
            showCancelButton: true,
            focusConfirm: false,
            willOpen: () => {
                if (willOpen) {
                    willOpen()
                }
            },
            preConfirm: () => {
                if (preConfirm) {
                    preConfirm()
                }
            },
        });
        if (results) {
            if (results.dismiss !== Swal.DismissReason.cancel) {
                if (results.value !== "") {
                    callback(results)
                } else {
                    callback(false)
                }
            } else {
                callback(false)
            }
        }
    };
    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    };
}