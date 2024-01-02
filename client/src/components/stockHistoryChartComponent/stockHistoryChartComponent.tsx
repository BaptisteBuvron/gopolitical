import React from "react";
import {Line} from "react-chartjs-2";
import {CategoryScale, Chart, Legend, LinearScale, LineElement, PointElement, Title, Tooltip,} from 'chart.js';
import {LabelColorResourceService} from "../../services/LabelColorResourceService";
Chart.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

interface StockHistoryChartProps {
    stockHistory: Map<number, Map<string, number>>,
}
function StockHistoryChart({ stockHistory }: StockHistoryChartProps) {
    const chartData = transformStockHistoryToChartData(stockHistory);


    const chartOptions = {
        responsive: true,
        plugins: {
            legend: {
                position: 'top' as const,
            },
        },
        scales: {
            x: {
                title: {
                    display: true,
                    text: 'Day',
                },
            },
            y: {
                title: {
                    display: true,
                    text: 'Resource Value',
                },
            },
        },
    };

    return <Line data={chartData} options={chartOptions} ></Line>;
}

function transformStockHistoryToChartData(stockHistory: Map<number, Map<string, number>>) {
    const labels: string[] = [];
    const datasets: any[] = [];

    // Extract labels and datasets from stockHistory
    stockHistory.forEach((stock, key) => {
        labels.push(key.toString()); // Assuming the keys are numeric, convert to string for labels

        // Iterate over each resource in the stock
        stock.forEach((value, resource) => {
            // Find or create a dataset for each resource
            const datasetIndex = datasets.findIndex(dataset => dataset.label === resource);

            if (datasetIndex !== -1) {
                // Dataset already exists, add data point
                datasets[datasetIndex].data.push(value);
            } else {
                // Dataset doesn't exist, create a new one
                datasets.push({
                    label: resource,
                    data: [value],
                    fill: false, // You can adjust these options based on your chart requirements
                    borderColor: LabelColorResourceService.getLabelColorByResource(resource),
                    pointRadius: 1,
                    pointHoverRadius: 7,
                });
            }
        });
    });

    return { labels, datasets };
}

export default StockHistoryChart;

