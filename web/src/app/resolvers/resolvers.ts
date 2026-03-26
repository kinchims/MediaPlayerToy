import { ActivatedRouteSnapshot, MaybeAsync, ResolveFn, RouterStateSnapshot } from "@angular/router"
import { of } from "rxjs";
import { Config } from "../models/config.model";
import { inject } from "@angular/core";
import { MPTAPI } from "../../services/api.service";

export class Resolvers {
    public static Config: ResolveFn<Config> = 
       (route: ActivatedRouteSnapshot, state: RouterStateSnapshot): MaybeAsync<Config> => {
                const api = inject(MPTAPI)
                console.log(api)
                return api.getSettings();
        }
}