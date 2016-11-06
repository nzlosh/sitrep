console.log(Date() + "[siterep.js] loading.");

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
    "get admins": "/api/admins",
    "get app version": "/api/version",
    "get all alerts": "/api/alertlog",
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

// Review Alerts
$("#load_alertlog").api({
    action: "get all alerts",
    method: "GET",
    onSuccess: function(response){
        max_pages = 10;
        page_size = response.length;
        console.log("Items: " + page_size);
        display_items = 0;

        start_exec = new Date();

        $("#review_table").empty();
        response.forEach(function(r){

            if ( display_items > parseInt(page_size/max_pages)) return;

            if ( r.status == "CRITICAL" ) {
                td_state = "class='error'";
                label_state = "red";
            } else if ( r.status == "WARNING" ) {
                td_state = "class='warning'";
                label_state = "orange";
            } else if ( r.status == "OK" ) {
                td_state = "class='positive'";
                label_state = "green";
            } else if ( r.status == "DOWN" ) {
                td_state = "class='error'";
                label_state = "pink";
            } else if ( r.status == "UP" ) {
                td_state = "class='positive'";
                label_state = "teal";
            } else {
                td_state = "";
                label_state = "grey";
            }
            $("#review_table").append("<tr id='"+r.id+"' " + td_state + "><td>" + new Date(r.alert_date*1000) + "</td>" +
                "<td>" + r.host + "</td>" +
                "<td>" + r.service + "</td>" +
                "<td><i class='ui "+ label_state + " label'>"  + r.status + "</i></td>" +
                "<td>" + r.output + "</td>" +
            "</tr>");
            display_items++;
        });

        console.log(parseInt(new Date() - start_exec));
    }
});


// Calendar for alerting period
$('#start_date').calendar({
  monthFirst: false,
  formatter: {
    date: function (date, settings) {
      if (!date) return '';
      var day = date.getDate();
      var month = date.getMonth() + 1;
      var year = date.getFullYear();
      return day + '/' + month + '/' + year;
    }
  }
});

$('#end_date').calendar({
  monthFirst: false,
  formatter: {
    date: function (date, settings) {
      if (!date) return '';
      var day = date.getDate();
      var month = date.getMonth() + 1;
      var year = date.getFullYear();
      return day + '/' + month + '/' + year;
    }
  }
});


// Application version from backend.
$("#app_version").api({
    action: "get app version",
    onSuccess: function(response){
        $("#app_version").empty();
        $("#app_version").append(r);
    }
});

console.log(Date() + "[siterep.js] loaded.");
