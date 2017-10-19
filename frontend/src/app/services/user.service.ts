import { Injectable } from "@angular/core";
import { Http } from "@angular/http";
import { Observable } from "rxjs/Observable";
import 'rxjs/add/operator/map';
import { User } from "../types/types";

const backend = 'http://dev:9000'

@Injectable()
export class UserService {
    constructor(private http: Http) {}

    getUsers(): Observable<User[]> {
        return this.http.get(`${backend}/usage`)
        .map(res => res.json() as User[])
    }
}