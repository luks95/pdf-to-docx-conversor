export namespace main {
	
	export class ConversionResult {
	    fileName: string;
	    success: boolean;
	    path: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ConversionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileName = source["fileName"];
	        this.success = source["success"];
	        this.path = source["path"];
	        this.error = source["error"];
	    }
	}

}

