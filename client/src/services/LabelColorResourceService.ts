export class LabelColorResourceService {
    static resourcesColor = [
        { resource: "petrol", color: "#ABCDEF" },
        { resource: "water", color: "#15A6CE" },
        { resource: "food", color: "#FF08B9" },
        { resource: "armement", color: "#8338E3" },
    ]
    static getLabelColorByResource(resourceName: string | undefined): string {
        const color = this.resourcesColor.find((resource) => resource.resource === resourceName)?.color
        return color ? color : "#8100B8";
    }
}