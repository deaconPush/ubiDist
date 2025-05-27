export namespace hdwallet {
	
	export class WalletTransaction {
	    sender: string;
	    recipient: string;
	    status: string;
	    value: string;
	    token: string;
	    createdAt: string;
	
	    static createFrom(source: any = {}) {
	        return new WalletTransaction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sender = source["sender"];
	        this.recipient = source["recipient"];
	        this.status = source["status"];
	        this.value = source["value"];
	        this.token = source["token"];
	        this.createdAt = source["createdAt"];
	    }
	}

}

