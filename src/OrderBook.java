import java.util.ArrayList;
import java.util.PriorityQueue;

public class OrderBook {
    private final PriorityQueue<Order> buys = new PriorityQueue<>();
    private final PriorityQueue<Order> sells = new PriorityQueue<>();
    private final ArrayList<Transaction> transactions = new ArrayList<>();
    private final String security;

    public OrderBook(String security) {
        this.security = security;
    }

    public void placeBuy(int amt, double price, String creator, int createdAt) {
        while (sells.size() > 0 && sells.peek().getPrice() <= price && amt > 0) {
            Order top = sells.peek();
            int amtMatched = Math.min(top.getAmt(), amt);
            amt -= amtMatched;
            if (top.reduceAmt(amtMatched) == 0) sells.remove();

            transactions.add(new Transaction(
                    top.getCreator(), creator, security,
                    amtMatched, price, System.currentTimeMillis()
            ));
        }
        if (amt != 0) buys.add(new Order(amt, price, creator, createdAt, true));
    }

    public void placeSell(int amt, double price, String creator, int createdAt) {
        while (buys.size() > 0 && buys.peek().getPrice() >= price && amt > 0) {
            Order top = buys.peek();
            int amtMatched = Math.min(top.getAmt(), amt);
            amt -= amtMatched;
            if (top.reduceAmt(amtMatched) == 0) buys.remove();

            transactions.add(new Transaction(
                    creator, top.getCreator(), security,
                    amtMatched, price, System.currentTimeMillis()
            ));
        }
        if (amt != 0) sells.add(new Order(amt, price, creator, createdAt, false));
    }

    public double getMarketPrice() {
        return transactions.get(transactions.size() - 1).price;
    }

    public ArrayList<Transaction> getTransactions() {
        return transactions;
    }
}
