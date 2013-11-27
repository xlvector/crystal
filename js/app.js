var app = {
    addDiscreteBarChart: function(data_path, container) {
        $.ajax({
            url: data_path,
            dataType: "json",
            success: function(data){
                nv.addGraph(function() {
                    var chart = nv.models.discreteBarChart()
                        .x(function(d) { return d.label })
                        .y(function(d) { return d.value })
                        .staggerLabels(true)
                        .tooltips(false)
                        .showValues(true);

                    
                    d3.select(container)
                        .datum(data)
                        .transition().duration(500)
                        .call(chart);

                    nv.utils.windowResize(chart.update);

                    return chart;
                });    
            },
        });
    },

    addStackedBarChart: function(data_path, container) {
        $.ajax({
            url: data_path,
            dataType: "json",
            success: function(data){
                nv.addGraph(function() {
                    var chart = nv.models.multiBarChart().stacked(true);
                    
                    d3.select(container)
                        .datum(data)
                        .transition().duration(500)
                        .call(chart);

                    nv.utils.windowResize(chart.update);

                    return chart;
                });    
            },
        });
    },

    addScatterChart: function(data_path, container) {
        $.ajax({
            url: data_path,
            dataType: "json",
            success: function(data){
                nv.addGraph(function() {
                    var chart = nv.models.scatterChart()
                        .showDistX(true)
                        .showDistY(true)
                        .color(d3.scale.category10().range());

                    chart.xAxis.tickFormat(d3.format('.02f'));
                    chart.yAxis.tickFormat(d3.format('.02f'));

                    d3.select(container)
                        .datum(data)
                        .transition().duration(500)
                        .call(chart);

                    nv.utils.windowResize(chart.update);

                    return chart;
                });    
            },
        });
    },

    addChartSelector: function(list_url, container, appvar){
        $.ajax({
            url: list_url,
            dataType: "json",
            success: function(data){
                chart_dict = {}
                for(var i = 0; i < data.length; i++){
                    chart = data[i];
                    $(container).append("<option value=\"" + chart.name + "\">" + chart.name + "</option>");
                    chart_dict[chart.name] = chart.charts;
                }

                $(container).change(
                    function(){
                        charts = chart_dict[$(this).val()];
                        appvar.addStackedBarChart("/data/adult/" + charts[0], "#chart1 svg");
                        appvar.addStackedBarChart("/data/adult/" + charts[1], "#chart2 svg");
                    }
                );
            }
        })
    }
}