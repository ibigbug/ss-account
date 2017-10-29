import { Injectable } from "@angular/core";
import { Http } from "@angular/http";
import { Observable } from "rxjs/Observable";
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';

import { User } from "../types/types";

import { environment } from "../../environments/environment";

const backend = environment.production ? `http://${location.host}` : environment.backend

@Injectable()
export class UserService {
    constructor(private http: Http) { }

    getUsers(): Observable<User[]> {
        return this.http.get(`${backend}/usage`)
            .map(res => res.json() as User[])
    }

    createUser(user: User): Observable<User> {
        let f = new FormData()
        f.append('username', user.username)
        f.append('backend', user.backend)

        return this.http.post(`${backend}/register`, f)
            .catch(err => Observable.throw(err.text()))
            .map(res => {
                return res.json() as User
            })
    }
}