console.log(Date() + "[siterep.js] loading... 44");

sitrep = {
    start_date: parseInt(new Date().getTime()/1000),
    end_date: parseInt(new Date().getTime()/1000) - 3600,
    review: {
        page_size: 10,
        updatePageSize: function(){
            console.log("Page size adjusted")
            sitrep.review.page_size = $('#alerts_per_page').dropdown('get value');
        },
        maxPages:  function(){
            return sitrep.review.data.length / sitrep.review.page_size;
        },
        data: {},
        updatePage: function(){
            display_items = 0;
            if (sitrep.review.page_size == 0 || sitrep.review.page_size > sitrep.review.data.length)
            {
                display_items = sitrep.review.data.length;
            }
            else
            {
                display_items =  sitrep.review.page_size;
            }
            console.log("Page size: " + display_items + " limit=" + sitrep.review.page_size * 5);

            start_exec = new Date();

            $("#review_table").empty();
            var item_html = "";
            var display_vars = {
                "CRITICAL": {
                    td_state: "class='error'",
                    label_state: "red"
                },
                "WARNING": {
                    td_state: "class='warning'",
                    label_state: "orange"
                },
                "OK": {
                    td_state: "class='positive'",
                    label_state: "green"
                },
                "DOWN": {
                    td_state: "class='error'",
                    label_state: "pink"
                },
                "UP": {
                    td_state: "class='positive'",
                    label_state: "teal"
                },
                "default": {
                    td_state: "",
                    label_state: "grey"
                }
            };

            for(var i = 0; i < display_items; i++)
            {
                var r = sitrep.review.data[i];
                if (r === undefined) break;
                var x = display_vars[r.status] || display_vars["default"];

                item_html += "<tr id='" + r.id + "' " + x.td_state + ">" +
                "<td>" + new Date(r.alert_date*1000) + "</td>" +
                "<td id=\"host_row\"><b>" + r.host + "</b>" +
                "/" + r.service + "</td>" +
                "<td id=\"status_row\"><i class='ui "+ x.label_state + " label'>"  + r.status + "</i></td>" +
                "<td id=\"output_row\">" + r.output + "</td>" +
                "</tr>\n"
            }

            $("#review_table").append(item_html);

            console.log("Render time:" + parseInt(new Date() - start_exec));
            delete start_exec;
            delete item_html;
        }
    },
    onReady: function() {

        $('.ui.dropdown').dropdown();

        $('.ui.menu .ui.dropdown').dropdown({
            on: 'hover'
        });

        // Populate admin user list.
        $('.ui.dropdown.item').api("query");

        // Main tabmenu item selection
        $('.ui.menu a.item').on('click', function() {
            $(this).addClass('active').siblings().removeClass('active');
        });

        heatmap();
        // weekly();  // This is broken and needs to be fixed.
        top_alerts();

        $('#all_alert_percentage').progress({
            duration : 77,
            total    : 200,
            text     : {
                active: '{value} of {total} done'
            }
        });

        // http://jsfiddle.net/urf6P/3/
        $("#host_alerts").change(function() {
            console.log("Try to hide stuff without " + $(this).val());
            $("#table_alerts td#host_row:contains('" + $(this).val() + "')").parent().show();
            $("#table_alerts td#host_row:not(:contains('" + $(this).val() + "'))").parent().hide();
            console.log("Column Service:" + $("#service_alerts").val());
        });

        // Activate the tabs on the page.
        $('.ui.pointing.menu .item').tab();
        // Review page number selection.
        $('#alerts_page').on("click", "a.item", function(){
            console.log("View page number " + $(this).value);
        });
        // Number of items to display per page
        $('#alerts_per_page .ui.search.dropdown').dropdown({
            onChange: function(value, text, choice){
                sitrep.review.page_size = value;
                sitrep.review.updatePage();
            }
        });

        $('.ui.sticky').sticky({
            context: '#review_pane'
        });

        // Attach a résumé to an approval request:
        $('#carlos_june_2018').popup({
            boundary: '.report .b'
        });
        console.log(Date() + "[siterep.js] Initialised.");
    }
}





// Semantic UI API map
$.fn.api.settings.api = {
    "get admins": "/api/admins",
    "get app version": "/api/version",
    "get all alerts": "/api/alertlog?s={s}&e={e}",
};

// Admins dropdown menu.
$('.ui.dropdown.item.admin_menu').api({
    action: "get admins",
    on: "ready",
    method: "GET",
    onSuccess: function(response){
        console.log("Get Admins called.")
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
    urlData: {
        s: sitrep.start_date,
        e: sitrep.end_date
    },
    onComplete:  function(response){
        console.log("alertlog complete");
    },
    onSuccess: function(response){
        sitrep.review.data = response;
        sitrep.review.updatePage();
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
      console.log("start_date: "+date.toLocaleString()+"\nDay: "+day+" Month: "+month+" Year: "+year);
      sitrep.start_date = parseInt(date.getTime()/1000);
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
      console.log("start_date: "+date.toLocaleString()+"\nDay: "+day+" Month: "+month+" Year: "+year);
      sitrep.end_date = parseInt(date.getTime()/1000);
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


// http://bl.ocks.org/juan-cb/faf62e91e3c70a99a306
function top_alerts(){
    var categories= ['','Accessories', 'Audiophile', 'Camera & Photo', 'Cell Phones', 'Computers','eBook Readers','Gadgets','GPS & Navigation','Home Audio','Office Electronics','Portable Audio','Portable Video','Security & Surveillance','Service','Television & Video','Car & Vehicle'];
    var dollars = [213,209,190,179,156,209,190,179,213,209,190,179,156,209,190,190];
    var colors = ['#0000b4','#0082ca','#0094ff','#0d4bcf','#0066AE','#074285','#00187B','#285964','#405F83','#416545','#4D7069','#6E9985','#7EBC89','#0283AF','#79BCBF','#99C19E'];

    var grid = d3.range(25).map(function(i){
        return {'x1':0,'y1':0,'x2':0,'y2':480};
    });

    var tickVals = grid.map(function(d,i){
        if (i > 0) {
            i = i * 10;
        } else if (i === 0) {
            i = 100;
        }
        return i;
    });

    var xscale =    d3.scale.linear()
                    .domain([10,250])
                    .range([0,722]);

    var yscale =    d3.scale.linear()
                    .domain([0,categories.length])
                    .range([0,480]);

    var colorScale =    d3.scale.quantize()
                        .domain([0,categories.length])
                        .range(colors);

    var canvas =    d3.select('#weekly_barchart')
                    .append('svg')
                    .attr({'width':900,'height':550});

    var grids = canvas.append('g')
                .attr('id','grid')
                .attr('transform','translate(150,10)')
                .selectAll('line')
                .data(grid)
                .enter()
                .append('line')
                .attr({
                    'x1':function(d,i){ return i*30; },
                    'y1':function(d){ return d.y1; },
                    'x2':function(d,i){ return i*30; },
                    'y2':function(d){ return d.y2; },
                })
                .style({'stroke':'#adadad','stroke-width':'1px'});

    var xAxis = d3.svg.axis()
                .orient('bottom')
                .scale(xscale)
                .tickValues(tickVals);

    var yAxis = d3.svg.axis()
                .orient('left')
                .scale(yscale)
                .tickSize(2)
                .tickFormat(function(d,i){ return categories[i]; })
                .tickValues(d3.range(17));

    var y_xis = canvas.append('g')
                .attr("transform", "translate(150,0)")
                .attr('id','yaxis')
                .call(yAxis);

    var x_xis = canvas.append('g')
                .attr("transform", "translate(150,480)")
                .attr('id','xaxis')
                .call(xAxis);

    var chart = canvas.append('g')
                .attr("transform", "translate(150,0)")
                .attr('id','bars')
                .selectAll('rect')
                .data(dollars)
                .enter()
                .append('rect')
                .attr('height',19)
                .attr({'x':0,'y':function(d,i){ return yscale(i)+19; }})
                .style('fill',function(d,i){ return colorScale(i); })
                .attr('width',function(d){ return 0; });


    var transit =   d3.select("svg").selectAll("rect")
                    .data(dollars)
                    .transition()
                    .duration(1000)
                    .attr("width", function(d) {return xscale(d); });

    var transitext =    d3.select('#bars')
                        .selectAll('text')
                        .data(dollars)
                        .enter()
                        .append('text')
                        .attr({'x':function(d) {return xscale(d)-200; },'y':function(d,i){ return yscale(i)+35; }})
                        .text(function(d){ return d+"$"; }).style({'fill':'#fff','font-size':'14px'});
}

/*
function weekly(){
    var margin = {
        top: 20,
        right: 20,
        bottom: 70,
        left: 40};
    width = 600 - margin.left - margin.right;
    height = 300 - margin.top - margin.bottom;

    // Parse the date / time
    var parseDate = d3.time.format("%Y-%m").parse;
    var x = d3.scale.ordinal().rangeRoundBands([0, width], .05);
    var y = d3.scale.linear().range([height, 0]);

    var xAxis = d3.svg.axis()
        .scale(x)
        .orient("bottom")
        .tickFormat(d3.time.format("%Y-%m"));

    var yAxis = d3.svg.axis()
        .scale(y)
        .orient("left")
        .ticks(10);

    var svg = d3.select("#all_alert_percentage").append("svg")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
        .append("g")
        .attr("transform",
              "translate(" + margin.left + "," + margin.top + ")");

    d3.csv("bar-data.csv", function(error, data) {
        data.date = parseDate(data.date);
        data.value =+ data.value;
    });

    x.domain(data.map(function(d) { return d.date; }));
    y.domain([0, d3.max(data, function(d) { return d.value; })]);

    svg.append("g")
        .attr("class", "x axis")
        .attr("transform", "translate(0," + height + ")")
        .call(xAxis)
        .selectAll("text")
        .style("text-anchor", "end")
        .attr("dx", "-.8em")
        .attr("dy", "-.55em")
        .attr("transform", "rotate(-90)" );

    svg.append("g")
        .attr("class", "y axis")
        .call(yAxis)
        .append("text")
        .attr("transform", "rotate(-90)")
        .attr("y", 6)
        .attr("dy", ".71em")
        .style("text-anchor", "end")
        .text("Value ($)");

    svg.selectAll("bar")
        .data(data)
        .enter().append("rect")
        .style("fill", "steelblue")
        .attr("x", function(d) { return x(d.date); })
        .attr("width", x.rangeBand())
        .attr("y", function(d) { return y(d.value); })
        .attr("height", function(d) { return height - y(d.value); });

}

*/

// Reference http://bl.ocks.org/tjdecke/5558084
function heatmap(){
    var margin = { top: 50, right: 0, bottom: 100, left: 50 },
    width = 960 - margin.left - margin.right,
    height = 430 - margin.top - margin.bottom,
    gridSize = Math.floor(width / 24),
    legendElementWidth = gridSize*2,
    buckets = 9,
    colors = ["#00a8ff", "#16b0ff", "#33b8ff", "#3299ff", "#0099ff", "#0066ff", "#0000ff", "#000099"],
    days = ["Mon", "Tue", "Wed", "Thur", "Fri", "Sat", "Sun"],
    times = ["00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23"];
    datasets = ["data.tsv", "data2.tsv"];

      var svg = d3.select("#heatmap_chart").append("svg")
          .attr("width", width + margin.left + margin.right)
          .attr("height", height + margin.top + margin.bottom)
          .append("g")
          .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

      var dayLabels = svg.selectAll(".dayLabel")
          .data(days)
          .enter().append("text")
            .text(function (d) { return d; })
            .attr("x", 0)
            .attr("y", function (d, i) { return i * gridSize; })
            .style("text-anchor", "end")
            .attr("transform", "translate(-6," + gridSize / 1.5 + ")")
            .attr("class", function (d, i) { return ((i >= 0 && i <= 4) ? "dayLabel mono axis axis-workweek" : "dayLabel mono axis"); });

      var timeLabels = svg.selectAll(".timeLabel")
          .data(times)
          .enter().append("text")
            .text(function(d) { return d; })
            .attr("x", function(d, i) { return i * gridSize; })
            .attr("y", 0)
            .style("text-anchor", "middle")
            .attr("transform", "translate(" + gridSize / 2 + ", -6)")
            .attr("class", function(d, i) { return ((i >= 7 && i <= 16) ? "timeLabel mono axis axis-worktime" : "timeLabel mono axis"); });

      var heatmapChart = function(tsvFile) {
        d3.tsv(tsvFile,
        function(d) {
          return {
            day: +d.day,
            hour: +d.hour,
            value: +d.value
          };
        },
        function(error, data) {
          var colorScale = d3.scale.quantile()
              .domain([0, buckets - 1, d3.max(data, function (d) { return d.value; })])
              .range(colors);

          var cards = svg.selectAll(".hour")
              .data(data, function(d) {return d.day+':'+d.hour;});

          cards.append("title");

          cards.enter().append("rect")
              .attr("x", function(d) { return (d.hour - 1) * gridSize; })
              .attr("y", function(d) { return (d.day - 1) * gridSize; })
              .attr("rx", 4)
              .attr("ry", 4)
              .attr("class", "hour bordered")
              .attr("width", gridSize)
              .attr("height", gridSize)
              .style("fill", colors[0]);

          cards.transition().duration(1000)
              .style("fill", function(d) { return colorScale(d.value); });

          cards.select("title").text(function(d) { return d.value; });

          cards.exit().remove();

          var legend = svg.selectAll(".legend")
              .data([0].concat(colorScale.quantiles()), function(d) { return d; });

          legend.enter().append("g")
              .attr("class", "legend");

          legend.append("rect")
            .attr("x", function(d, i) { return legendElementWidth * i; })
            .attr("y", height)
            .attr("width", legendElementWidth)
            .attr("height", gridSize / 2)
            .style("fill", function(d, i) { return colors[i]; });

          legend.append("text")
            .attr("class", "mono")
            .text(function(d) { return "≥ " + Math.round(d); })
            .attr("x", function(d, i) { return legendElementWidth * i; })
            .attr("y", height + gridSize);
          legend.exit().remove();
        });
      };
      heatmapChart(datasets[0]);
};

console.log(Date() + "[siterep.js] loaded.");
$(document).ready(function() {
    // Initialise webapp
    sitrep.onReady();
});
