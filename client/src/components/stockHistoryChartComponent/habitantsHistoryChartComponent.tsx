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

interface HabitantsHistoryChartProps {
    habitantsHistory: Map<string, number>,
}
function StockHistoryChart({ habitantsHistory }: HabitantsHistoryChartProps) {



    const chartData = transformHabitantsHistoryToChartData(habitantsHistory);


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
                    text: 'Habitants Number',
                },
            },
        },
    };

    return <Line data={chartData} options={chartOptions} ></Line>;
}


function transformHabitantsHistoryToChartData(habitantsHistory: Map<string, number>) {
    const labels: string[] = [];
    const datasets: any[] = [];

    // Extract labels and datasets from habitantsHistory
    for (let [key] of habitantsHistory) {
        labels.push(key); // Assuming the keys are numeric, convert to string for labels
    }

    datasets.push({
        label: "Habitants",
        data: Array.from(habitantsHistory.values()),
        fill: false, // You can adjust these options based on your chart requirements
        borderColor: LabelColorResourceService.getLabelColorByResource("Habitants"),
        pointRadius: 1,
        pointHoverRadius: 7,
    });

    return { labels, datasets };
}

export default StockHistoryChart;

