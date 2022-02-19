public class Transaction {
    final String from;
    final String to;
    final String security;
    final double price;
    final int amt;
    final long time;

    public Transaction(String from, String to, String security, int amt, double price, long time) {
        this.from = from;
        this.to = to;
        this.security = security;
        this.amt = amt;
        this.price = price;
        this.time = time;
    }

    public String toString() {
        return from + " -> " + to + ": " + amt + " " + security + " @ $" + price + " at time " + time;
    }
}
