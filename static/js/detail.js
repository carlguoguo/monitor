var apiID=$("#api-ID").val()
$.getJSON('/data/api/request/'+apiID, function (data) {
    // $.each(data, function(index) {
    //     dateString = data[index][0]
    //     dateTimeParts = dateString.split(' '),
    //     timeParts = dateTimeParts[1].split(':'),
    //     dateParts = dateTimeParts[0].split('-'),
    //     date = new Date(dateParts[0], dateParts[1], dateParts[2], timeParts[0], timeParts[1], timeParts[2]);
    //     data[index][0] = date
    // });
    Highcharts.chart('container', {
        chart: {
            zoomType: 'x'
        },
        title: {
            text: '响应时间'
        },
        subtitle: {
            text: '拖拽可以放大聚焦'
        },
        xAxis: {
            visible: false
        },
        yAxis: {
            title: {
                text: '响应时间(毫秒)'
            },
        },
        legend: {
            enabled: false
        },
        plotOptions: {
            area: {
                fillColor: {
                    linearGradient: {
                        x1: 0,
                        y1: 0,
                        x2: 0,
                        y2: 1
                    },
                    stops: [
                        [0, Highcharts.getOptions().colors[0]],
                        [1, Highcharts.Color(Highcharts.getOptions().colors[0]).setOpacity(0).get('rgba')]
                    ]
                },
                marker: {
                    radius: 2
                },
                lineWidth: 1,
                states: {
                    hover: {
                        lineWidth: 1
                    }
                },
                threshold: null
            }
        },

        series: [{
            type: 'area',
            name: '响应时间(毫秒):',
            data: data,
        }]
    });
});
