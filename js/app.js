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

    selected_dataset : "",

    addDataSetSelector: function(dataset_url, list_url, dataset_container, feature_container, appvar){
        $.ajax({
            url : dataset_url, 
            dataType: "json",
            success: function(data) {
                $(dataset_container).append("<option value=\"none\"></option>");
                for(var i = 0; i < data.length; i++){
                    dataset = data[i];
                    $(dataset_container).append("<option value=\"" + dataset + "\">" + dataset + "</option>");
                }

                $(dataset_container).change(
                    function(){
                        selected_dataset = $(this).val();
                        if(selected_dataset != "none"){
                            appvar.selected_dataset = selected_dataset;
                            appvar.addChartSelector(list_url, selected_dataset, feature_container, appvar);
                        }
                    }
                );
            }
        })
    },

    addChartSelector: function(list_url, selected_dataset, container, appvar){
        $(container).html("");
        $.ajax({
            url: list_url,
            data: {dataset: selected_dataset},
            dataType: "json",
            success: function(data){
                for(var i = 0; i < data.length; i++){
                    chart = data[i];
                    $(container).append("<option value=\"" + chart + "\">" + chart + "</option>");
                    if(i == 0){
                        appvar.addStackedBarChart("/single_feature?dataset=" + selected_dataset 
                            + "&feature=" + chart, "#chart1 svg");
                        appvar.addStackedBarChart("/single_feature?dataset=" + selected_dataset 
                            + "&is_100percent=true&feature=" + chart, "#chart2 svg");
                    }
                }

                $(container).change(
                    function(){
                        chart = $(this).val();
                        appvar.addStackedBarChart("/single_feature?dataset=" + selected_dataset 
                            + "&feature=" + chart, "#chart1 svg");
                        appvar.addStackedBarChart("/single_feature?dataset=" + selected_dataset 
                            + "&is_100percent=true&feature=" + chart, "#chart2 svg");
                    }
                );
            }
        })
    }
}