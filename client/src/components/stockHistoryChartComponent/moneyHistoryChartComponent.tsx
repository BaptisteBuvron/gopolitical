import React from "react";
import { Line } from "react-chartjs-2";
import {
    CategoryScale,
    Chart,
    Legend,
    LinearScale,
    LineElement,
    PointElement,
    Title,
    Tooltip,
} from "chart.js";

Chart.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

interface MoneyHistoryChartProps {
    moneyHistory: Map<string, number>;
}

function MoneyHistoryChart({ moneyHistory }: MoneyHistoryChartProps) {
    const chartData = transformMoneyHistoryToChartData(moneyHistory);

    const chartOptions = {
        responsive: true,
        plugins: {
            legend: {
                position: "top" as const,
            },
        },
        scales: {
            x: {
                title: {
                    display: true,
                    text: "Day",
                },
            },
            y: {
                title: {
                    display: true,
                    text: "Resource Value",
                },
            },
        },
    };

    return <Line data={chartData} options={chartOptions}></Line>;
}

function transformMoneyHistoryToChartData(moneyHistory: Map<string, number>) {
    const labels: string[] = Array.from(moneyHistory.keys());
    const datasets: any[] = [
        {
            label: "Money",
            data: Array.from(moneyHistory.values()),
            fill: false,
            borderColor: "blue",
            pointRadius: 1,
            pointHoverRadius: 7,
        },
    ];

    return { labels, datasets };
}




export default MoneyHistoryChart;
