public class Order implements Comparable<Order> {
    private final boolean buy;
    private final String creator;
    private final int createdAt;
    private int amt;
    private final double price;

    public Order(int amt, double price, String creator, int createdAt, boolean buy) {
        this.price = price;
        this.amt = amt;
        this.creator = creator;
        this.createdAt = createdAt;
        this.buy = buy;
    }

    public int compareTo(Order o) {
        if (o.price != price) {
            if (buy) {
                if (price < o.price) return 1;
            } else if (price > o.price) return 1;
            return -1;
        }

        if (createdAt < o.createdAt) return -1;
        else if (createdAt > o.createdAt) return 1;

        return 0;
    }

    public int getAmt() {
        return amt;
    }

    public double getPrice() {
        return price;
    }

    public int reduceAmt(int amt) {
        this.amt -= amt;
        return this.amt;
    }

    public String getCreator() {
        return creator;
    }
}
