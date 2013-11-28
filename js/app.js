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
        $(container).html("")
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
                for(var i = 0; i < data.length; i++){
                    chart = data[i];
                    $(container).append("<option value=\"" + chart + "\">" + chart + "</option>");
                }

                $(container).change(
                    function(){
                        chart = $(this).val();
                        appvar.addStackedBarChart("/single_feature?feature=" + chart, "#chart1 svg");
                        appvar.addStackedBarChart("/single_feature?feature=" + chart, "#chart2 svg");
                    }
                );
            }
        })
    }
}