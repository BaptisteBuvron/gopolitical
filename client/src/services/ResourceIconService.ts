import {resourceIcons} from "../resourceIcons";

export class ResourceIconService {
    private resourceIcons: {resource: string, iconPath: string}[];

    constructor() {
        this.resourceIcons = resourceIcons;
    }

    getResourceIconPath(resource: string): string {
        const resourceIcon = this.resourceIcons.find((r) => r.resource === resource);
        return resourceIcon ? resourceIcon.iconPath : "";
    }
}