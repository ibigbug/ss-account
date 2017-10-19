import { Http } from "@angular/http";
import { Injectable } from "@angular/core";

@Injectable()
export class AppConfig {
    private config: Object = null;
    private env: Object = null;
    constructor(private http: Http) {}

    public getEnv(key: any) {
        return this.env[key]
    }

    public load() {}
}