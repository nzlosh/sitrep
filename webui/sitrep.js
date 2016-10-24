console.log("[siterep.js] loading.");

// Activate dropdown menu.
$(document).ready(function() {
    $('.ui.dropdown').dropdown();

    $('.ui.menu .ui.dropdown').dropdown({
        on: 'hover'
    });

    $('.ui.menu a.item').on('click', function() {
        $(this).addClass('active').siblings().removeClass('active');
    });
});

// Activate the tables on the page.
$('.ui.pointing.menu .item').tab();

// Semantic UI API map
$.fn.api.settings.api = {
    'get admins': '/api/admins'
};

// Admins dropdown menu.
$('.ui.dropdown.item').api({
    action: "get admins",
    method: "GET",
    onSuccess: function(response){
        $("#admin_menu").empty();
        response.forEach(function(r){
            if (r.is_active){
                $("#admin_menu").append('<div class="item">' + r.login + '</div>');
            }
        });
    }
});

$('#start_date').calendar();
$('#end_date').calendar();

$("#app_version").api({
    action: "get app version",
    onSuccess: function(response){
        $("#app_version").empty();
        $("#app_version").append(r);
    }
});

console.log("[siterep.js] loaded.");
